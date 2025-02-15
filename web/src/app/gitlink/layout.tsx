import HomeNav from '@/components/HomeNav';
import React from 'react';


export default function PublicLayout({ children, modal }: React.PropsWithChildren<{
    modal: React.ReactNode;
  }>) {
    return <>
        <HomeNav initialKey="gitlink"/>
        {children}
    </>
  }