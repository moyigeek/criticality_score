import clsx from 'clsx';
import styles from './layout.module.css';
import React from 'react';

export default function PublicLayout({ children }: React.PropsWithChildren<{}>) {
  return <div className={clsx(styles['public-container'])}>
    {children}
  </div>
}