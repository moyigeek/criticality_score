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
  // 你可以添加更多的菜单项
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