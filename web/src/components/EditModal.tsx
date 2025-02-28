"use client";
import React, { useState } from "react";
import { Modal, Input, Card,Skeleton, Form, Row, Col } from "antd";
import { postUpdateGitlink } from "@/service/client";
import { Data } from "./GitLinkTable"; // 导入 Data 类型


interface EditModalProps {
    currentPackage: Data;
    tableName: string;
    onClose: () => void;
}

const EditModal: React.FC<EditModalProps> = ({ currentPackage, tableName, onClose }) => {
    const [visible, setVisible] = useState(true);
    const [loading, setLoading] = useState(false);
    const [form] = Form.useForm();

    const handleOk = async () => {
        try {
            const values = await form.validateFields();
            setLoading(true);
            await postUpdateGitlink({
                body: {
                    newGitLink: values.gitlink,
                    packageName: currentPackage.package,
                    tableName: tableName,
                    linkConfidence: values.link_confidence, 
                }
            });
            setLoading(false);
            setVisible(false);
            form.resetFields();
            onClose();
        } catch (error) {
            setLoading(false);
            console.error("Failed to update gitlink:", error);
        }
    };

    const handleCancel = () => {
        setVisible(false);
        form.resetFields();
        onClose();
    };

    return (
        <Modal
            open={visible}
            title="Edit GitLink"
            onOk={handleOk}
            onCancel={handleCancel}
            confirmLoading={loading}
            okText="Confirm"
            cancelText="Cancel"
            width='70%'
            height='50%'
        >
            <Row gutter={16}>
                <Col span={12}>
                    <div>
                        <p><strong>Package Name:</strong> {currentPackage.package}</p>
                        <p><strong>Description:</strong> {currentPackage.description}</p>
                        <p><strong>Homepage:</strong> <a href={currentPackage.homepage} target="_blank" rel="noopener noreferrer">{currentPackage.homepage}</a></p>
                        <p><strong>Current GitLink:</strong> <a href={currentPackage.git_link} target="_blank" rel="noopener noreferrer">{currentPackage.git_link}</a></p>
                        <p><strong>Table Name:</strong> {tableName}</p>
                    </div>
                </Col>
                <Col span={12}>
                    <Row>
                        <Card title="AI advice"  style={{ width: 300 }}>
                        <Skeleton active />
                        </Card>
                    </Row>
                    <Row>
                    <Form form={form} layout="vertical" name="edit_gitlink_form">
                        <Form.Item
                            name="gitlink"
                            label="New GitLink"
                            rules={[{ required: true, message: "Please input the new gitlink!" }]}
                        >
                            <Input />
                        </Form.Item>
                        <Form.Item
                            name="link_confidence"
                            label="Link Confidence"
                            rules={[{ required: true, message: "Please input the link confidence!" }]}
                            initialValue={1}
                        >
                            <Input />
                        </Form.Item>
                    </Form>
                    </Row>
                </Col>
            </Row>
        </Modal>
    );
};

export default EditModal;