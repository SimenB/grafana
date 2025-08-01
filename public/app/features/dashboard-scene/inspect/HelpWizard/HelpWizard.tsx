import { css } from '@emotion/css';
import { useMemo, useEffect } from 'react';
import AutoSizer from 'react-virtualized-auto-sizer';

import { GrafanaTheme2, FeatureState } from '@grafana/data';
import { Trans, t } from '@grafana/i18n';
import { config } from '@grafana/runtime';
import { VizPanel } from '@grafana/scenes';
import {
  Drawer,
  Tab,
  TabsBar,
  CodeEditor,
  useStyles2,
  Field,
  InlineSwitch,
  Button,
  Spinner,
  Alert,
  FeatureBadge,
  Select,
  ClipboardButton,
  Stack,
  TextLink,
} from '@grafana/ui';
import { contextSrv } from 'app/core/services/context_srv';
import { AccessControlAction } from 'app/types/accessControl';

import { ShowMessage, SnapshotTab, SupportSnapshotService } from './SupportSnapshotService';

interface Props {
  panel: VizPanel;
  onClose: () => void;
}

export function HelpWizard({ panel, onClose }: Props) {
  const styles = useStyles2(getStyles);
  const service = useMemo(() => new SupportSnapshotService(panel), [panel]);
  const plugin = panel.getPlugin();

  const {
    currentTab,
    loading,
    error,
    options,
    showMessage,
    snapshotSize,
    markdownText,
    snapshotText,
    randomize,
    panelTitle,
    scene,
  } = service.useState();

  useEffect(() => {
    service.buildDebugDashboard();
  }, [service, plugin, randomize]);

  if (!plugin) {
    return null;
  }

  const tabs = [
    { label: t('dashboard-scene.help-wizard.tabs.label.snapshot', 'Snapshot'), value: SnapshotTab.Support },
    { label: t('dashboard-scene.help-wizard.tabs.label.data', 'Data'), value: SnapshotTab.Data },
  ];

  const hasSupportBundleAccess =
    config.supportBundlesEnabled && contextSrv.hasPermission(AccessControlAction.ActionSupportBundlesCreate);

  return (
    <Drawer
      title={t('dashboard-scene.help-wizard.title-get-help-with-this-panel', 'Get help with this panel')}
      size="lg"
      onClose={onClose}
      subtitle={
        <Stack direction="column" gap={1}>
          <Stack direction="row" gap={1}>
            <FeatureBadge featureState={FeatureState.beta} />
            <TextLink href="https://grafana.com/docs/grafana/latest/troubleshooting/" external>
              <Trans i18nKey="dashboard-scene.help-wizard.troubleshooting-docs">Troubleshooting docs</Trans>
            </TextLink>
          </Stack>
          <span className="muted">
            <Trans i18nKey="dashboard-scene.help-wizard.troubleshooting-request-help">
              To request troubleshooting help, send a snapshot of this panel to Grafana Labs Technical Support. The
              snapshot contains query response data and panel settings.
            </Trans>
          </span>
          {hasSupportBundleAccess && (
            <span className="muted">
              <Trans i18nKey="dashboard-scene-help-wizard.support-bundle">
                You can also retrieve a support bundle containing information concerning your Grafana instance and
                configured datasources in the <TextLink href="/support-bundles">support bundles section</TextLink>.
              </Trans>
            </span>
          )}
        </Stack>
      }
      tabs={
        <TabsBar>
          {tabs.map((t, index) => (
            <Tab
              key={`${t.value}-${index}`}
              label={t.label}
              active={t.value === currentTab}
              onChangeTab={() => service.onCurrentTabChange(t.value!)}
            />
          ))}
        </TabsBar>
      }
    >
      {loading && <Spinner />}
      {error && <Alert title={error.title}>{error.message}</Alert>}

      {currentTab === SnapshotTab.Data && (
        <div className={styles.code}>
          <div className={styles.opts}>
            <Field label={t('dashboard-scene.help-wizard.label-template', 'Template')} className={styles.field}>
              <Select options={options} value={showMessage} onChange={service.onShowMessageChange} />
            </Field>

            {showMessage === ShowMessage.GithubComment ? (
              <ClipboardButton icon="copy" getText={service.onGetMarkdownForClipboard}>
                <Trans i18nKey="dashboard-scene.help-wizard.copy-to-clipboard">Copy to clipboard</Trans>
              </ClipboardButton>
            ) : (
              <Button icon="download-alt" onClick={service.onDownloadDashboard}>
                <Trans i18nKey="dashboard-scene.help-wizard.download-snapshot">Download ({{ snapshotSize }})</Trans>
              </Button>
            )}
          </div>
          <AutoSizer disableWidth>
            {({ height }) => (
              <CodeEditor
                width="100%"
                height={height}
                language={showMessage === ShowMessage.GithubComment ? 'markdown' : 'json'}
                showLineNumbers={true}
                showMiniMap={true}
                value={showMessage === ShowMessage.GithubComment ? markdownText : snapshotText}
                readOnly={false}
                onBlur={service.onSetSnapshotText}
              />
            )}
          </AutoSizer>
        </div>
      )}
      {currentTab === SnapshotTab.Support && (
        <>
          <Field
            label={t('dashboard-scene.help-wizard.label-randomize-data', 'Randomize data')}
            description={t(
              'dashboard-scene.help-wizard.description-randomize-data',
              'Modify the original data to hide sensitive information. Note the lengths will stay the same, and duplicate values will be equal.'
            )}
          >
            <Stack>
              <InlineSwitch
                label={t('dashboard-scene.help-wizard.randomize-labels-label-labels', 'Labels')}
                id="randomize-labels"
                showLabel={true}
                value={Boolean(randomize.labels)}
                onChange={() => service.onToggleRandomize('labels')}
              />
              <InlineSwitch
                label={t('dashboard-scene.help-wizard.randomize-field-names-label-field-names', 'Field names')}
                id="randomize-field-names"
                showLabel={true}
                value={Boolean(randomize.names)}
                onChange={() => service.onToggleRandomize('names')}
              />
              <InlineSwitch
                label={t('dashboard-scene.help-wizard.randomize-string-values-label-string-values', 'String values')}
                id="randomize-string-values"
                showLabel={true}
                value={Boolean(randomize.values)}
                onChange={() => service.onToggleRandomize('values')}
              />
            </Stack>
          </Field>

          <Field
            label={t('dashboard-scene.help-wizard.label-support-snapshot', 'Support snapshot')}
            description={t('dashboard-scene.help-wizard.description-support-snapshot', 'Panel: {{panelTitle}}', {
              panelTitle,
            })}
          >
            <Stack>
              <Button icon="download-alt" onClick={service.onDownloadDashboard}>
                <Trans i18nKey="dashboard-scene.help-wizard.download-dashboard">Dashboard ({{ snapshotSize }})</Trans>
              </Button>
              <ClipboardButton
                icon="github"
                getText={service.onGetMarkdownForClipboard}
                title={t(
                  'dashboard-scene.help-wizard.title-complete-git-hub-comment-clipboard',
                  'Copy a complete GitHub comment to the clipboard'
                )}
              >
                <Trans i18nKey="dashboard-scene.help-wizard.copy-to-clipboard">Copy to clipboard</Trans>
              </ClipboardButton>
            </Stack>
          </Field>

          <div style={{ height: '100%', overflow: 'auto', display: 'flex' }}>
            {scene && <scene.Component model={scene} />}
          </div>
        </>
      )}
    </Drawer>
  );
}

const getStyles = (theme: GrafanaTheme2) => {
  return {
    code: css({
      flexGrow: 1,
      height: '100%',
      overflow: 'scroll',
    }),
    field: css({
      width: '100%',
    }),
    opts: css({
      display: 'flex',
      width: '100%',
      flexGrow: 0,
      alignItems: 'center',
      justifyContent: 'flex-end',

      '& button': {
        marginLeft: theme.spacing(1),
      },
    }),
  };
};
