import clsx from 'clsx';
import styles from './layout.module.css';
import React from 'react';
import TopNav from '@/components/TopNav';

export default function PublicLayout({ children, modal }: React.PropsWithChildren<{
  modal: React.ReactNode;
}>) {
  return <>
    <TopNav initialKey="home"/>
    <div className={clsx(styles['public-container'])}>
      {children}
    </div>
    {modal}
  </>
}