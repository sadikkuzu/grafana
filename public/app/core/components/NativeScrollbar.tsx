import { css, cx } from '@emotion/css';
import { useEffect, useRef } from 'react';

import { CustomScrollbar, useStyles2 } from '@grafana/ui';

type Props = Parameters<typeof CustomScrollbar>[0];

// Shim to provide API-compatibility for Page's scroll-related props
export default function NativeScrollbar({ children, scrollRefCallback, scrollTop, divId }: Props) {
  const styles = useStyles2(getStyles);
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (ref.current && scrollRefCallback) {
      scrollRefCallback(ref.current);
    }
  }, [ref, scrollRefCallback]);

  useEffect(() => {
    if (ref.current && scrollTop != null) {
      ref.current?.scrollTo(0, scrollTop);
    }
  }, [scrollTop]);

  return (
    // Set the .scrollbar-view class to help e2e tests find this, like in CustomScrollbar
    <div ref={ref} className={cx(styles.nativeScrollbars, 'scrollbar-view')} id={divId}>
      {children}
    </div>
  );
}

function getStyles() {
  return {
    nativeScrollbars: css({
      label: 'native-scroll-container',
      minHeight: `calc(100% + 0px)`, // I don't know, just copied from custom scrollbars
      maxHeight: `calc(100% + 0px)`, // I don't know, just copied from custom scrollbars
      display: 'flex',
      flexDirection: 'column',
      flexGrow: 1,
      overflow: 'auto',
    }),
  };
}
