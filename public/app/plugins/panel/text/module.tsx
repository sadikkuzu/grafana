import { PanelPlugin } from '@grafana/data';

import { TextPanel } from './TextPanel';
import { TextPanelEditor } from './TextPanelEditor';
import { defaultPanelOptions, PanelOptions, TextMode } from './models.gen';
import { textPanelMigrationHandler } from './textPanelMigrationHandler';

export const plugin = new PanelPlugin<PanelOptions>(TextPanel)
  .setPanelOptions((builder) => {
    builder
      .addRadio({
        path: 'mode',
        name: 'Mode',
        description: 'text mode of the panel',
        settings: {
          options: [
            { value: TextMode.Markdown, label: 'Markdown' },
            { value: TextMode.HTML, label: 'HTML' },
          ],
        },
        defaultValue: defaultPanelOptions.mode,
      })
      .addCustomEditor({
        id: 'content',
        path: 'content',
        name: 'Content',
        description: 'Content of the panel',
        editor: TextPanelEditor,
        defaultValue: defaultPanelOptions.content,
      })
      .addBooleanSwitch({
        path: 'query',
        name: 'Execute query',
        description: 'Run a query and include the results in template variable context',
        defaultValue: defaultPanelOptions.query,
      });
  })
  .setSkipDataQuery((opts?: PanelOptions) => {
    return opts?.query !== true;
  })
  .setMigrationHandler(textPanelMigrationHandler);
