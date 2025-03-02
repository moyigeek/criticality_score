"use client";
import { getQueryWithPagination } from "@/service/client";
import { Table, Switch, Button, Row, Col } from "antd";
import React, { useEffect, useState } from "react";
import TableDropdown from "./TableDropdown";
import EditModal from "./EditModal";

type Data = {
    package: string,
    description: string,
    homepage: string,
    git_link: string,
    link_confidence: string,
    key: string, // 添加 key 属性
}

const GitLinkTable: React.FC = () => {
    const [data, setData] = useState<Data[]>([]);
    const [currentPage, setCurrentPage] = useState<number>(1);
    const [pageSize, setPageSize] = useState<number>(10);
    const [total, setTotal] = useState<number>(0);
    const [tableName, setTableName] = useState<string>("arch_packages");
    const [selectedPackage, setSelectedPackage] = useState<Data | null>(null);
    const [confidence, setConfidence] = useState<boolean>(true);

    const fetchData = async (page: number, pageSize: number) => {
        const response = await getQueryWithPagination({
            query: {
                tableName: tableName, // 替换为实际的表名
                pageSize: pageSize,
                offset: (page - 1) * pageSize,
                confidence: confidence,
            },
        });
        if (response && response.data && Array.isArray(response.data.items)) {
            const itemsWithKey = response.data.items.map((item: Data, index: number) => ({
                ...item,
                key: item.package, // 使用 package 作为 key
            }));
            setData(itemsWithKey);
            setTotal(response.data.totalPages as number * pageSize); // 假设 totalPages 是总页数
        }
    };

    useEffect(() => {
        fetchData(currentPage, pageSize);
    }, [currentPage, pageSize, tableName, confidence]);

    const handleTableChange = (pagination: any) => {
        setCurrentPage(pagination.current);
        setPageSize(pagination.pageSize);
    };

    const handleTableNameChange = (newTableName: string) => {
        setTableName(newTableName);
    };

    const handleEditClick = (record: Data) => {
        setSelectedPackage(record);
    };

    const handleModalClose = () => {
        setSelectedPackage(null);
        // 更新数据
        fetchData(currentPage, pageSize);
    };

    const handleConfidenceChange = (checked: boolean) => {
        setConfidence(checked);
    };

    return (
        <div>
           <Row gutter={[16, 16]} justify="space-between" align="middle">
                <Col>
                    <TableDropdown onTableNameChange={handleTableNameChange} tableName={tableName} />
                </Col>
                <Col>
                    <Switch
                        checkedChildren="linkConfidence"
                        unCheckedChildren="linkConfidence"
                        checked={confidence}
                        onChange={handleConfidenceChange}
                    />
                </Col>
            </Row>
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
                    { title: 'Description', dataIndex: 'description', key: 'description', width: '20%' },
                    { title: 'Homepage', dataIndex: 'homepage', key: 'homepage' },
                    { title: 'Git Link', dataIndex: 'git_link', key: 'git_link', width: '25%' },
                    { title: 'Link Confidence', dataIndex: 'link_confidence', key: 'link_confidence', width: '10%' },
                    {
                        title: 'Action',
                        key: 'action',
                        render: (_, record) => (
                            <Button type="link" onClick={() => handleEditClick(record)}>
                                Edit
                            </Button>
                        ),
                    },
                ]}
            />
            {selectedPackage && (
                <EditModal
                    currentPackage={selectedPackage}
                    tableName={tableName}
                    onClose={handleModalClose}
                />
            )}
        </div>
    );
};

export default GitLinkTable;