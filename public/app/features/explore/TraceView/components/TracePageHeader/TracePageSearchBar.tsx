// Copyright (c) 2018 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import { css } from '@emotion/css';
import cx from 'classnames';
import React, { memo, Dispatch, SetStateAction } from 'react';

import { GrafanaTheme2 } from '@grafana/data';
import { config, reportInteraction } from '@grafana/runtime';
import { Button, useStyles2 } from '@grafana/ui';

import { ubJustifyEnd } from '../uberUtilityStyles';

export type TracePageSearchBarProps = {
  // searchValue: string;
  // setSearch: (value: string) => void;
  spanFindMatches: Set<string> | undefined;
  focusedSpanIdForSearch: string;
  setFocusedSpanIdForSearch: Dispatch<SetStateAction<string>>;
  datasourceType?: string;
};

export default memo(function TracePageSearchBar(props: TracePageSearchBarProps) {
  const {
    // searchValue,
    spanFindMatches,
    focusedSpanIdForSearch,
    setFocusedSpanIdForSearch,
    datasourceType,
  } = props;
  const styles = useStyles2(getStyles);

  // const btnClass = cx(styles.TracePageSearchBarBtn, { [styles.TracePageSearchBarBtnDisabled]: !searchValue });
  const btnClass = cx(styles.TracePageSearchBarBtn);

  // const setTraceSearch = (value: string) => {
  //   setFocusedSpanIdForSearch('');
  // };

  const nextResult = () => {
    reportInteraction('grafana_traces_trace_view_find_next_prev_clicked', {
      datasourceType: datasourceType,
      grafana_version: config.buildInfo.version,
      direction: 'next',
    });

    const spanMatches = Array.from(spanFindMatches!);
    const prevMatchedIndex = spanMatches.indexOf(focusedSpanIdForSearch)
      ? spanMatches.indexOf(focusedSpanIdForSearch)
      : 0;

    // new query || at end, go to start
    if (prevMatchedIndex === -1 || prevMatchedIndex === spanMatches.length - 1) {
      setFocusedSpanIdForSearch(spanMatches[0]);
      return;
    }

    // get next
    setFocusedSpanIdForSearch(spanMatches[prevMatchedIndex + 1]);
  };

  const prevResult = () => {
    reportInteraction('grafana_traces_trace_view_find_next_prev_clicked', {
      datasourceType: datasourceType,
      grafana_version: config.buildInfo.version,
      direction: 'prev',
    });

    const spanMatches = Array.from(spanFindMatches!);
    const prevMatchedIndex = spanMatches.indexOf(focusedSpanIdForSearch)
      ? spanMatches.indexOf(focusedSpanIdForSearch)
      : 0;

    // new query || at start, go to end
    if (prevMatchedIndex === -1 || prevMatchedIndex === 0) {
      setFocusedSpanIdForSearch(spanMatches[spanMatches.length - 1]);
      return;
    }

    // get prev
    setFocusedSpanIdForSearch(spanMatches[prevMatchedIndex - 1]);
  };

  return (
    <div className={styles.TracePageSearchBar}>
      <span className={ubJustifyEnd} style={{ display: 'flex' }}>
        <>
          <Button
            className={btnClass}
            variant="secondary"
            // disabled={!searchValue}
            disabled={false}
            type="button"
            icon="arrow-down"
            aria-label="Next results button"
            onClick={nextResult}
          />
          <Button
            className={btnClass}
            variant="secondary"
            // disabled={!searchValue}
            disabled={false}
            type="button"
            icon="arrow-up"
            aria-label="Prev results button"
            onClick={prevResult}
          />
        </>
      </span>
    </div>
  );
});

export const getStyles = (theme: GrafanaTheme2) => {
  return {
    TracePageSearchBar: css`
      label: TracePageSearchBar;
      float: right;
      position: absolute;
      top: 0;
      right: 0;
      z-index: ${theme.zIndex.navbarFixed};
      background: ${theme.colors.background.primary};
      margin-bottom: -48px;
      padding: 8px;
      margin-right: 2px;
      border-radius: ${theme.shape.borderRadius()};
      box-shadow: ${theme.shadows.z2};
    `,
    TracePageSearchBarBar: css`
      label: TracePageSearchBarBar;
      max-width: 20rem;
      transition: max-width 0.5s;
      &:focus-within {
        max-width: 100%;
      }
    `,
    TracePageSearchBarBtn: css`
      label: TracePageSearchBarBtn;
      transition: 0.2s;
      margin-left: 8px;
    `,
    TracePageSearchBarBtnDisabled: css`
      label: TracePageSearchBarBtnDisabled;
      opacity: 0.5;
    `,
    TracePageSearchBarLocateBtn: css`
      label: TracePageSearchBarLocateBtn;
      padding: 1px 8px 4px;
    `,
  };
};
