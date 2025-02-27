"use client";
import { getQueryWithPagination } from "@/service/client";
import { Table } from "antd";
import React, { useEffect, useState } from "react";

type Data = {
    package: string,
    description: string,
    homepage: string,
    git_link: string,
    key: string, // 添加 key 属性
}

const GitLinkTable: React.FC = () => {
    const [data, setData] = useState<Data[]>([]);
    const [currentPage, setCurrentPage] = useState<number>(1);
    const [pageSize, setPageSize] = useState<number>(10);
    const [total, setTotal] = useState<number>(0);
    const [tableName, setTableName] = useState<string>("arch_packages");

    const fetchData = async (page: number, pageSize: number) => {
        const response = await getQueryWithPagination({
            query: {
                tableName: tableName, // 替换为实际的表名
                pageSize: pageSize,
                offset: (page - 1) * pageSize,
            },
        });
        if (response && response.data && Array.isArray(response.data.items)) {
            const itemsWithKey = response.data.items.map((item: Data, index: number) => ({
                ...item,
                key: item.package, // 使用 package 作为 key
            }));
            setData(itemsWithKey);
            setTotal(response.data.totalPages as number); // 假设 totalPages 是总页数
        }
    };

    useEffect(() => {
        fetchData(currentPage, pageSize);
    }, [currentPage, pageSize]);

    const handleTableChange = (pagination: any) => {
        setCurrentPage(pagination.current);
        setPageSize(pagination.pageSize);
    };

    return (
        <div>
            <Table
                dataSource={data}
                pagination={{
                    current: currentPage,
                    pageSize: pageSize,
                    total: total,
                    showQuickJumper: true,
                }}
                onChange={handleTableChange}
                columns={[
                    { title: 'Package', dataIndex: 'package', key: 'package' },
                    { title: 'Description', dataIndex: 'description', key: 'description' },
                    { title: 'Homepage', dataIndex: 'homepage', key: 'homepage' },
                    { title: 'Git Link', dataIndex: 'git_link', key: 'git_link' },
                    { title: 'Action', key: 'action', render: () => <a>View</a> },
                ]}
            />
        </div>
    );
};

export default GitLinkTable;