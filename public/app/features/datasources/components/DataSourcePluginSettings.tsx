import { createElement, PureComponent } from 'react';

import { DataSourcePluginMeta, DataSourceSettings } from '@grafana/data';
import { writableProxy } from 'app/features/plugins/extensions/utils';

import { GenericDataSourcePlugin } from '../types';

export interface Props {
  plugin: GenericDataSourcePlugin;
  dataSource: DataSourceSettings;
  dataSourceMeta: DataSourcePluginMeta;
  onModelChange: (dataSource: DataSourceSettings) => void;
}

export class DataSourcePluginSettings extends PureComponent<Props> {
  constructor(props: Props) {
    super(props);

    this.onModelChanged = this.onModelChanged.bind(this);
  }

  onModelChanged = (dataSource: DataSourceSettings) => {
    this.props.onModelChange(dataSource);
  };

  render() {
    const { plugin, dataSource } = this.props;

    if (!plugin) {
      return null;
    }

    return (
      <div>
        {plugin.components.ConfigEditor &&
          createElement(plugin.components.ConfigEditor, {
            options: writableProxy(dataSource, { source: 'datasource', pluginId: plugin.meta?.id }),
            onOptionsChange: this.onModelChanged,
          })}
      </div>
    );
  }
}
