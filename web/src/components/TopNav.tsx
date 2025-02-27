"use client";
import React, { useState, useEffect } from 'react';
import type { MenuProps } from 'antd';
import { Menu } from 'antd';
import Link from 'next/link';

type MenuItem = Required<MenuProps>['items'][number];

const items: MenuItem[] = [
  {
    label: <Link href="/">Home</Link>,
    key: 'home',
  },
  {
    label: <Link href="/gitlink">GitLink</Link>,
    key: 'gitlink',
  },
];

interface HomeNavProps {
  initialKey?: string;
}

const TopNav: React.FC<HomeNavProps> = ({ initialKey = 'home' }) => {
  const [current, setCurrent] = useState(initialKey);

  useEffect(() => {
    setCurrent(initialKey);
  }, [initialKey]);

  const onClick: MenuProps['onClick'] = (e) => {
    setCurrent(e.key);
  };

  return <Menu onClick={onClick} selectedKeys={[current]} mode="horizontal" items={items} />;
};

export default TopNav;