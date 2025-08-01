package cloudmigrationimpl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/grafana/alerting/definition"
	"github.com/prometheus/alertmanager/pkg/labels"
	"github.com/stretchr/testify/require"

	"github.com/grafana/grafana/pkg/apimachinery/errutil"
	"github.com/grafana/grafana/pkg/components/simplejson"
	"github.com/grafana/grafana/pkg/services/accesscontrol"
	"github.com/grafana/grafana/pkg/services/dashboards"
	"github.com/grafana/grafana/pkg/services/datasources"
	"github.com/grafana/grafana/pkg/services/featuremgmt"
	"github.com/grafana/grafana/pkg/services/folder"
	ac "github.com/grafana/grafana/pkg/services/ngalert/accesscontrol"
	"github.com/grafana/grafana/pkg/services/ngalert/api/tooling/definitions"
	"github.com/grafana/grafana/pkg/services/ngalert/models"
	"github.com/grafana/grafana/pkg/services/user"
	"github.com/grafana/grafana/pkg/setting"
)

// Read-only.
var alertRulesPermissions = map[string][]string{
	accesscontrol.ActionAlertingRuleRead:   {"*"},
	accesscontrol.ActionAlertingRuleCreate: {"*"},
	accesscontrol.ActionAlertingRuleUpdate: {"*"},
	dashboards.ActionFoldersRead:           {"*"},
	datasources.ActionQuery:                {"*"},
}

func TestGetAlertMuteTimings(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("it returns the mute timings", func(t *testing.T) {
		t.Parallel()

		s := setUpServiceTest(t, false).(*Service)
		s.features = featuremgmt.WithFeatures(featuremgmt.FlagOnPremToCloudMigrations)

		user := &user.SignedInUser{OrgID: 1}

		createdMuteTiming := createMuteTiming(t, ctx, s, user)

		muteTimeIntervals, err := s.getAlertMuteTimings(ctx, user)
		require.NoError(t, err)
		require.NotNil(t, muteTimeIntervals)
		require.Len(t, muteTimeIntervals, 1)
		require.Equal(t, createdMuteTiming.Name, muteTimeIntervals[0].Name)
	})
}

func TestGetNotificationTemplates(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("it returns the notification templates", func(t *testing.T) {
		t.Parallel()

		s := setUpServiceTest(t, false).(*Service)

		user := &user.SignedInUser{OrgID: 1}

		createdTemplate := createNotificationTemplate(t, ctx, s, user)

		notificationTemplates, err := s.getNotificationTemplates(ctx, user)
		require.NoError(t, err)
		require.NotNil(t, notificationTemplates)
		require.Len(t, notificationTemplates, 1)
		require.Equal(t, createdTemplate.Name, notificationTemplates[0].Name)
	})
}

func TestGetContactPoints(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("it returns the contact points", func(t *testing.T) {
		t.Parallel()

		s := setUpServiceTest(t, false).(*Service)

		user := &user.SignedInUser{
			OrgID: 1,
			Permissions: map[int64]map[string][]string{
				1: {
					accesscontrol.ActionAlertingNotificationsRead:    nil,
					accesscontrol.ActionAlertingReceiversReadSecrets: {ac.ScopeReceiversAll},
				},
			},
		}

		createdContactPoints := createContactPoints(t, ctx, s, user)

		contactPoints, err := s.getContactPoints(ctx, user)
		require.NoError(t, err)
		require.NotNil(t, contactPoints)
		require.Len(t, contactPoints, len(createdContactPoints))
	})

	t.Run("it returns an error when user lacks permission to read contact point secrets", func(t *testing.T) {
		t.Parallel()

		s := setUpServiceTest(t, false).(*Service)

		user := &user.SignedInUser{
			OrgID: 1,
			Permissions: map[int64]map[string][]string{
				1: {
					accesscontrol.ActionAlertingNotificationsRead: nil,
				},
			},
		}

		contactPoints, err := s.getContactPoints(ctx, user)
		require.Nil(t, contactPoints)

		gfErr := errutil.Error{}
		require.ErrorAs(t, err, &gfErr)
		require.Equal(t, http.StatusForbidden, gfErr.Reason.Status().HTTPStatus())
	})
}

func TestGetNotificationPolicies(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("it returns the contact points", func(t *testing.T) {
		t.Parallel()

		s := setUpServiceTest(t, false).(*Service)

		user := &user.SignedInUser{OrgID: 1}

		muteTiming := createMuteTiming(t, ctx, s, user)
		require.NotEmpty(t, muteTiming.Name)

		contactPoints := createContactPoints(t, ctx, s, user)
		require.GreaterOrEqual(t, len(contactPoints), 1)

		updateNotificationPolicyTree(t, ctx, s, user, contactPoints[0].Name, muteTiming.Name)

		notificationPolicies, err := s.getNotificationPolicies(ctx, user)
		require.NoError(t, err)
		require.NotEmpty(t, notificationPolicies.Routes.Receiver)
		require.NotNil(t, notificationPolicies.Routes.Routes)
	})
}

func TestGetAlertRules(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("it returns the alert rules", func(t *testing.T) {
		t.Parallel()

		s := setUpServiceTest(t, false).(*Service)

		user := &user.SignedInUser{OrgID: 1, Permissions: map[int64]map[string][]string{1: alertRulesPermissions}}

		alertRule := createAlertRule(t, ctx, s, user, false, "")

		alertRules, err := s.getAlertRules(ctx, user)
		require.NoError(t, err)
		require.Len(t, alertRules, 1)
		require.Equal(t, alertRule.UID, alertRules[0].UID)
	})

	t.Run("when the alert_rules_state config is `paused`, then the alert rules are all returned in `paused` state", func(t *testing.T) {
		t.Parallel()

		alertRulesState := func(c *setting.Cfg) {
			c.CloudMigration.AlertRulesState = setting.GMSAlertRulesPaused
		}

		s := setUpServiceTest(t, false, alertRulesState).(*Service)

		user := &user.SignedInUser{OrgID: 1, Permissions: map[int64]map[string][]string{1: alertRulesPermissions}}

		alertRulePaused := createAlertRule(t, ctx, s, user, true, "")
		require.True(t, alertRulePaused.IsPaused)

		alertRuleUnpaused := createAlertRule(t, ctx, s, user, false, "")
		require.False(t, alertRuleUnpaused.IsPaused)

		alertRules, err := s.getAlertRules(ctx, user)
		require.NoError(t, err)
		require.Len(t, alertRules, 2)
		require.True(t, alertRules[0].IsPaused)
		require.True(t, alertRules[1].IsPaused)
	})
}

func TestGetAlertRuleGroups(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("it returns the alert rule groups", func(t *testing.T) {
		t.Parallel()

		s := setUpServiceTest(t, false).(*Service)

		user := &user.SignedInUser{OrgID: 1, Permissions: map[int64]map[string][]string{1: alertRulesPermissions}}

		ruleGroupTitle := "ruleGroupTitle"

		alertRule1 := createAlertRule(t, ctx, s, user, true, ruleGroupTitle)
		alertRule2 := createAlertRule(t, ctx, s, user, false, ruleGroupTitle)
		alertRule3 := createAlertRule(t, ctx, s, user, false, "anotherRuleGroup")

		createAlertRuleGroup(t, ctx, s, user, ruleGroupTitle, []models.AlertRule{alertRule1, alertRule2})

		ruleGroups, err := s.getAlertRuleGroups(ctx, user)
		require.NoError(t, err)
		require.Len(t, ruleGroups, 2)

		for _, ruleGroup := range ruleGroups {
			alertRuleUIDs := make([]string, 0)
			for _, alertRule := range ruleGroup.Rules {
				alertRuleUIDs = append(alertRuleUIDs, alertRule.UID)
			}

			if ruleGroup.Title == ruleGroupTitle {
				require.Len(t, ruleGroup.Rules, 2)
				require.ElementsMatch(t, []string{alertRule1.UID, alertRule2.UID}, alertRuleUIDs)
			} else {
				require.Len(t, ruleGroup.Rules, 1)
				require.ElementsMatch(t, []string{alertRule3.UID}, alertRuleUIDs)
			}
		}
	})

	t.Run("with the alert rules state set to paused, it returns the alert rule groups with alert rules paused", func(t *testing.T) {
		t.Parallel()

		alertRulesState := func(c *setting.Cfg) {
			c.CloudMigration.AlertRulesState = setting.GMSAlertRulesPaused
		}

		s := setUpServiceTest(t, false, alertRulesState).(*Service)

		user := &user.SignedInUser{OrgID: 1, Permissions: map[int64]map[string][]string{1: alertRulesPermissions}}

		ruleGroupTitle := "ruleGroupTitle"

		alertRule1 := createAlertRule(t, ctx, s, user, true, ruleGroupTitle)
		alertRule2 := createAlertRule(t, ctx, s, user, false, ruleGroupTitle)
		alertRule3 := createAlertRule(t, ctx, s, user, false, "anotherRuleGroup")

		createAlertRuleGroup(t, ctx, s, user, ruleGroupTitle, []models.AlertRule{alertRule1, alertRule2})

		ruleGroups, err := s.getAlertRuleGroups(ctx, user)
		require.NoError(t, err)
		require.Len(t, ruleGroups, 2)

		for _, ruleGroup := range ruleGroups {
			alertRuleUIDs := make([]string, 0)
			for _, alertRule := range ruleGroup.Rules {
				alertRuleUIDs = append(alertRuleUIDs, alertRule.UID)

				require.True(t, alertRule.IsPaused)
			}

			if ruleGroup.Title == ruleGroupTitle {
				require.Len(t, ruleGroup.Rules, 2)
				require.ElementsMatch(t, []string{alertRule1.UID, alertRule2.UID}, alertRuleUIDs)
			} else {
				require.Len(t, ruleGroup.Rules, 1)
				require.ElementsMatch(t, []string{alertRule3.UID}, alertRuleUIDs)
			}
		}
	})
}

func createMuteTiming(t *testing.T, ctx context.Context, service *Service, user *user.SignedInUser) definitions.MuteTimeInterval {
	t.Helper()

	muteTiming := `{
		"name": "My Unique MuteTiming 1",
		"time_intervals": [
			{
				"times": [{"start_time": "12:12","end_time": "23:23"}],
				"weekdays": ["monday","wednesday","friday","sunday"],
				"days_of_month": ["10:20","25:-1"],
				"months": ["1:6","10:12"],
				"years": ["2022:2054"],
				"location": "Africa/Douala"
			}
		]
	}`

	var mt definitions.MuteTimeInterval
	require.NoError(t, json.Unmarshal([]byte(muteTiming), &mt))

	createdTiming, err := service.ngAlert.Api.MuteTimings.CreateMuteTiming(ctx, mt, user.GetOrgID())
	require.NoError(t, err)

	return createdTiming
}

func createNotificationTemplate(t *testing.T, ctx context.Context, service *Service, user *user.SignedInUser) definitions.NotificationTemplate {
	t.Helper()

	tmpl := definitions.NotificationTemplate{
		Name:     "MyTestNotificationTemplate",
		Template: "This is a test template\n{{ .ExternalURL }}",
	}

	createdTemplate, err := service.ngAlert.Api.Templates.CreateTemplate(ctx, user.GetOrgID(), tmpl)
	require.NoError(t, err)

	return createdTemplate
}

func createContactPoints(t *testing.T, ctx context.Context, service *Service, user *user.SignedInUser) []definitions.EmbeddedContactPoint {
	t.Helper()

	slackSettings, err := simplejson.NewJson([]byte(`{
		"icon_emoji":"iconemoji",
		"icon_url":"iconurl",
		"recipient":"recipient",
		"token":"slack-secret",
		"username":"user"
	}`))
	require.NoError(t, err)

	telegramSettings, err := simplejson.NewJson([]byte(`{
		"bottoken":"telegram-secret",
		"chatid":"chat-id",
		"disable_notification":true,
		"disable_web_page_preview":false,
		"message_thread_id":"1234",
		"parse_mode":"None",
		"protect_content":true
	}`))
	require.NoError(t, err)

	nameGroup := "group_1"

	slackContactPoint := definitions.EmbeddedContactPoint{
		Name:                  nameGroup,
		Type:                  "slack",
		Settings:              slackSettings,
		DisableResolveMessage: false,
		Provenance:            "",
	}

	createdSlack, err := service.ngAlert.Api.ContactPointService.CreateContactPoint(ctx, user.GetOrgID(), user, slackContactPoint, "")
	require.NoError(t, err)

	telegramContactPoint := definitions.EmbeddedContactPoint{
		Name:                  nameGroup,
		Type:                  "telegram",
		Settings:              telegramSettings,
		DisableResolveMessage: false,
		Provenance:            "",
	}

	createdTelegram, err := service.ngAlert.Api.ContactPointService.CreateContactPoint(ctx, user.GetOrgID(), user, telegramContactPoint, "")
	require.NoError(t, err)

	return []definitions.EmbeddedContactPoint{
		createdSlack,
		createdTelegram,
	}
}

func updateNotificationPolicyTree(t *testing.T, ctx context.Context, service *Service, user *user.SignedInUser, receiverGroup, muteTiming string) {
	t.Helper()

	child := definition.Route{
		Continue:          true,
		MuteTimeIntervals: []string{muteTiming},
		ObjectMatchers: definition.ObjectMatchers{
			{Name: "label1", Type: labels.MatchEqual, Value: "value1"},
			{Name: "label2", Type: labels.MatchNotEqual, Value: "value2"},
		},
		Receiver: receiverGroup,
	}

	tree := definition.Route{
		Receiver: "grafana-default-email",
		Routes:   []*definition.Route{&child},
	}

	_, _, err := service.ngAlert.Api.Policies.UpdatePolicyTree(ctx, user.GetOrgID(), tree, "", "")
	require.NoError(t, err)
}

func createAlertRule(t *testing.T, ctx context.Context, service *Service, user *user.SignedInUser, isPaused bool, ruleGroup string) models.AlertRule {
	t.Helper()

	// Ensure the folder exists before creating alert rules
	createFolder(t, ctx, service, user, "folderUID", "Test Folder")

	rule := models.AlertRule{
		OrgID:        user.GetOrgID(),
		Title:        fmt.Sprintf("Alert Rule SLO (Paused: %v) - %v", isPaused, ruleGroup),
		NamespaceUID: "folderUID",
		Condition:    "A",
		Data: []models.AlertQuery{
			{
				RefID: "A",
				Model: []byte(`{"queryType": "a"}`),
				RelativeTimeRange: models.RelativeTimeRange{
					From: models.Duration(60),
					To:   models.Duration(0),
				},
			},
		},
		IsPaused:        isPaused,
		RuleGroup:       ruleGroup,
		For:             time.Minute,
		IntervalSeconds: 60,
		NoDataState:     models.OK,
		ExecErrState:    models.OkErrState,
	}

	createdRule, err := service.ngAlert.Api.AlertRules.CreateAlertRule(ctx, user, rule, "")
	require.NoError(t, err)

	return createdRule
}

func createFolder(t *testing.T, ctx context.Context, service *Service, user *user.SignedInUser, uid, title string) {
	t.Helper()
	_, err := service.folderService.Create(ctx, &folder.CreateFolderCommand{
		OrgID:        user.GetOrgID(),
		UID:          uid,
		Title:        title,
		SignedInUser: user,
	})
	if err != nil && !errors.Is(err, dashboards.ErrFolderWithSameUIDExists) {
		require.NoError(t, err)
	}
}

func createAlertRuleGroup(t *testing.T, ctx context.Context, service *Service, user *user.SignedInUser, title string, rules []models.AlertRule) models.AlertRuleGroup {
	t.Helper()

	group := models.AlertRuleGroup{
		Title:     title,
		FolderUID: "folderUID",
		Interval:  300,
		Rules:     rules,
	}

	err := service.ngAlert.Api.AlertRules.ReplaceRuleGroup(ctx, user, group, "")
	require.NoError(t, err)

	return group
}
