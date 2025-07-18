import { useCallback, useMemo } from 'react';
import { useEffectOnce } from 'react-use';

import { SelectableValue } from '@grafana/data';
import { t } from '@grafana/i18n';
import { Select } from '@grafana/ui';

import { selectors } from '../../e2e/selectors';
import { ResultFormat } from '../../types/query';
import { FormatAsFieldProps } from '../../types/types';

import { Field } from './Field';

const FormatAsField = ({
  query,
  variableOptionGroup,
  onQueryChange,
  inputId,
  options: formatOptions,
  defaultValue,
  onLoad,
  setFormatAs,
  resultFormat,
}: FormatAsFieldProps) => {
  const options = useMemo(() => [...formatOptions, variableOptionGroup], [variableOptionGroup, formatOptions]);

  const handleChange = useCallback(
    (change: SelectableValue<ResultFormat>) => {
      const { value } = change;
      if (!value) {
        return;
      }

      const newQuery = setFormatAs(query, value);
      onQueryChange(newQuery);
    },
    [onQueryChange, query, setFormatAs]
  );

  useEffectOnce(() => {
    //sets to default if the value is not found in the list
    if (!formatOptions.find((item) => item.value === resultFormat)) {
      handleChange({ value: defaultValue });
    }
    onLoad(query, defaultValue, handleChange);
  });

  return (
    <Field
      label={t('components.format-as-field.label-format-as', 'Format as')}
      data-testid={selectors.components.queryEditor.logsQueryEditor.formatSelection.input}
    >
      <Select
        inputId={`${inputId}-format-as-field`}
        value={resultFormat}
        onChange={handleChange}
        options={options}
        width={20}
      />
    </Field>
  );
};

export default FormatAsField;
