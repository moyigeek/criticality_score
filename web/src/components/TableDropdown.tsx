import React from 'react';
import { DownOutlined } from '@ant-design/icons';
import type { MenuProps } from 'antd';
import { Dropdown, Space, Menu } from 'antd';
import { table } from 'console';

const items: MenuProps['items'] = [
  {
    key: '1',
    label: 'alpine_packages',
  },
  {
    key: '2',
    label: 'arch_packages',
  },
  {
    key: '3',
    label: 'aur_packages',
  },
  {
    key: '4',
    label: 'centos_packages',
  },
  {
    key: '5',
    label: 'debian_packages',
  },
  {
    key: '6',
    label: 'deepin_packages',
  },
  {
    key: '7',
    label: 'fedora_packages',
  },
  {
    key: '8',
    label: 'gentoo_packages',
  },
  {
    key: '9',
    label: 'git_packages',
  },
  {
    key: '10',
    label: 'homebrew_packages',
  },
  {
    key:'11',
    label:'nix_packages',
  },
  {
    key:'12',
    label:'ubuntu_packages',
  }
  
];

interface TableDropdownProps {
  tableName: string;
  onTableNameChange: (tableName: string) => void;
}

const TableDropdown: React.FC<TableDropdownProps> = ({ tableName, onTableNameChange }) => {
  const handleMenuClick: MenuProps['onClick'] = (e) => {
    const selectedItem = items.find(item => item?.key === e.key);
    if (selectedItem && 'label' in selectedItem) {
      onTableNameChange(selectedItem.label as string);
    }
  };

  return (
    <Dropdown menu={{ items, onClick: handleMenuClick }}>
      <a onClick={(e) => e.preventDefault()}>
        <Space>
          {tableName}
          <DownOutlined />
        </Space>
      </a>
    </Dropdown>
  );
};

export default TableDropdown;