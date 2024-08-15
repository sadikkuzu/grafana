import React from 'react';
import { testIds } from '../components/testIds';
import { PluginPage, usePluginComponent } from '@grafana/runtime';

type ReusableComponentProps = {
  name: string;
};

export function ExposedComponents() {
  var { component: ReusableComponent } = usePluginComponent<ReusableComponentProps>(
    'myorg-c-app/reusable-component/v1'
  );

  if (!ReusableComponent) {
    return null;
  }

  return (
    <PluginPage>
      <div data-testid={testIds.pageTwo.container}>
        <ReusableComponent name={'World'} />
      </div>
    </PluginPage>
  );
}
