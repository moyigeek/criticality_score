import React from 'react';
import { DownOutlined, SmileOutlined } from '@ant-design/icons';
import type { MenuProps } from 'antd';
import { Dropdown, Space } from 'antd';

const items: MenuProps['items'] = [
  {
    key: '1',
    label: 'Arch Linux',
  },
  
];

const DistroSel: React.FC = () => (
  <Dropdown menu={{ items }}>
    <a onClick={(e) => e.preventDefault()}>
      <Space>
        Select table
        <DownOutlined />
      </Space>
    </a>
  </Dropdown>
);

export default DistroSel;