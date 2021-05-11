import React from 'react';
import { DataFrame, TimeRange } from '@grafana/data';
import { GraphNG, GraphNGProps } from '../GraphNG/GraphNG';
import { UPlotConfigBuilder } from '../uPlot/config/UPlotConfigBuilder';
import { PlotLegend } from '../uPlot/PlotLegend';
import { LegendDisplayMode } from '../VizLegend/models.gen';
import { preparePlotConfigBuilder } from './utils';
import { withTheme2 } from '../../themes/ThemeContext';
import { PanelContext, PanelContextRoot } from '../PanelChrome/PanelContext';

const propsToDiff: string[] = [];

type TimeSeriesProps = Omit<GraphNGProps, 'prepConfig' | 'propsToDiff' | 'renderLegend'>;

export class UnthemedTimeSeries extends React.Component<TimeSeriesProps> {
  static contextType = PanelContextRoot;
  panelContext: PanelContext = {} as PanelContext;

  prepConfig = (alignedFrame: DataFrame, getTimeRange: () => TimeRange) => {
    const { eventBus, sync } = this.context;
    const { theme } = this.props;
    return preparePlotConfigBuilder({
      frame: alignedFrame,
      theme,
      timeZone: this.panelContext.timeZone,
      getTimeRange,
      eventBus,
      sync,
    });
  };

  renderLegend = (config: UPlotConfigBuilder) => {
    const { legend, onLegendClick, frames } = this.props;

    if (!config || (legend && legend.displayMode === LegendDisplayMode.Hidden)) {
      return;
    }

    return (
      <PlotLegend
        data={frames}
        config={config}
        onLegendClick={onLegendClick}
        maxHeight="35%"
        maxWidth="60%"
        {...legend}
      />
    );
  };

  render() {
    return (
      <GraphNG
        {...this.props}
        prepConfig={this.prepConfig}
        propsToDiff={propsToDiff}
        renderLegend={this.renderLegend as any}
      />
    );
  }
}

export const TimeSeries = withTheme2(UnthemedTimeSeries);
TimeSeries.displayName = 'TimeSeries';
