import React, { useEffect, useState } from 'react';
import {
  Table,
  Button,
  Modal,
  Form,
  Input,
  Select,
  message,
  Popconfirm,
  Space,
  Tag,
  Card,
} from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import { configApi } from '@/api/config';
import type { ConfigTemplate, TemplateVariable } from '@/types';

const ConfigTemplateList: React.FC = () => {
  const [templates, setTemplates] = useState<ConfigTemplate[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingTemplate, setEditingTemplate] = useState<ConfigTemplate | null>(null);
  const [form] = Form.useForm();

  const fetchTemplates = async () => {
    setLoading(true);
    try {
      const data = await configApi.listTemplates();
      setTemplates(data);
    } catch (error) {
      message.error('Failed to fetch config templates');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchTemplates();
  }, []);

  const handleCreate = () => {
    setEditingTemplate(null);
    form.resetFields();
    form.setFieldValue('variables', [
      {
        name: '',
        type: 'string',
        default: '',
        required: false,
        description: '',
        secret: false,
      },
    ]);
    setModalVisible(true);
  };

  const handleEdit = (record: ConfigTemplate) => {
    setEditingTemplate(record);
    form.setFieldsValue(record);
    setModalVisible(true);
  };

  const handleDelete = async (id: string) => {
    try {
      await configApi.deleteTemplate(id);
      message.success('Config template deleted successfully');
      fetchTemplates();
    } catch (error) {
      message.error('Failed to delete config template');
    }
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      if (editingTemplate) {
        await configApi.updateTemplate(editingTemplate.id, values);
        message.success('Config template updated successfully');
      } else {
        await configApi.createTemplate(values);
        message.success('Config template created successfully');
      }
      setModalVisible(false);
      fetchTemplates();
    } catch (error) {
      message.error('Failed to save config template');
    }
  };

  const columns: ColumnsType<ConfigTemplate> = [
    {
      title: 'Name',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: 'Description',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
    },
    {
      title: 'Adapter Type',
      dataIndex: 'adapter_type',
      key: 'adapter_type',
      render: (type: string) => <Tag color="blue">{type}</Tag>,
    },
    {
      title: 'Version',
      dataIndex: 'version',
      key: 'version',
    },
    {
      title: 'Variables',
      dataIndex: 'variables',
      key: 'variables',
      render: (variables: TemplateVariable[]) => (
        <span>{variables.length} variables</span>
      ),
    },
    {
      title: 'Actions',
      key: 'actions',
      render: (_, record) => (
        <Space>
          <Button
            type="text"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          >
            Edit
          </Button>
          <Popconfirm
            title="Are you sure you want to delete this config template?"
            onConfirm={() => handleDelete(record.id)}
            okText="Yes"
            cancelText="No"
          >
            <Button type="text" danger icon={<DeleteOutlined />}>
              Delete
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <Card title="Config Templates">
      <div style={{ marginBottom: 16 }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
          Create Template
        </Button>
      </div>
      <Table
        columns={columns}
        dataSource={templates}
        loading={loading}
        rowKey="id"
        pagination={{
          showSizeChanger: true,
          showQuickJumper: true,
        }}
      />
      <Modal
        title={editingTemplate ? 'Edit Config Template' : 'Create Config Template'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
        width={800}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="name"
            label="Name"
            rules={[{ required: true, message: 'Please input template name' }]}
          >
            <Input placeholder="Template name" />
          </Form.Item>
          <Form.Item name="description" label="Description">
            <Input.TextArea rows={3} placeholder="Template description" />
          </Form.Item>
          <Form.Item
            name="adapter_type"
            label="Adapter Type"
            rules={[{ required: true, message: 'Please select adapter type' }]}
          >
            <Select placeholder="Select adapter type">
              <Select.Option value="OpenClaw">OpenClaw</Select.Option>
              <Select.Option value="NanoClaw">NanoClaw</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name="version" label="Version">
            <Input placeholder="1.0.0" />
          </Form.Item>
          <Form.Item label="Variables">
            <Form.List name="variables">
              {(fields, { add, remove }) => (
                <>
                  {fields.map(({ key, name, ...restField }) => (
                    <Space key={key} style={{ display: 'flex', marginBottom: 8 }} align="baseline">
                      <Form.Item
                        {...restField}
                        name={[name, 'name']}
                        rules={[{ required: true, message: 'Required' }]}
                      >
                        <Input placeholder="Variable name" style={{ width: 150 }} />
                      </Form.Item>
                      <Form.Item
                        {...restField}
                        name={[name, 'type']}
                        rules={[{ required: true, message: 'Required' }]}
                      >
                        <Select placeholder="Type" style={{ width: 100 }}>
                          <Select.Option value="string">String</Select.Option>
                          <Select.Option value="number">Number</Select.Option>
                          <Select.Option value="boolean">Boolean</Select.Option>
                        </Select>
                      </Form.Item>
                      <Form.Item {...restField} name={[name, 'default']}>
                        <Input placeholder="Default" style={{ width: 120 }} />
                      </Form.Item>
                      <Form.Item {...restField} name={[name, 'description']} initialValue="">
                        <Input placeholder="Description" style={{ width: 150 }} />
                      </Form.Item>
                      <Form.Item {...restField} name={[name, 'secret']} valuePropName="checked" initialValue={false}>
                        <input type="checkbox" />
                      </Form.Item>
                      <Form.Item {...restField} name={[name, 'required']} valuePropName="checked" initialValue={false}>
                        <input type="checkbox" />
                      </Form.Item>
                      {fields.length > 1 && (
                        <Button
                          type="link"
                          danger
                          onClick={() => remove(name)}
                          icon={<DeleteOutlined />}
                        />
                      )}
                    </Space>
                  ))}
                  <Form.Item>
                    <Button type="dashed" onClick={() => add()} block icon={<PlusOutlined />}>
                      Add Variable
                    </Button>
                  </Form.Item>
                </>
              )}
            </Form.List>
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  );
};

export default ConfigTemplateList;