import React, { useState, useEffect } from 'react';
import {
  Table,
  Button,
  Modal,
  Form,
  Input,
  InputNumber,
  message,
  Popconfirm,
  Space,
  Card,
  Tabs,
} from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import { tenantApi, projectApi } from '@/api/tenant';
import type { Tenant, Project } from '@/types';

const TenantProjectManagement: React.FC = () => {
  const [activeTab, setActiveTab] = useState('tenants');
  const [tenants, setTenants] = useState<Tenant[]>([]);
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingItem, setEditingItem] = useState<Tenant | Project | null>(null);
  const [form] = Form.useForm();

  const fetchTenants = async () => {
    setLoading(true);
    try {
      const data = await tenantApi.listTenants();
      setTenants(data);
    } catch (error) {
      message.error('Failed to fetch tenants');
    } finally {
      setLoading(false);
    }
  };

  const fetchProjects = async () => {
    setLoading(true);
    try {
      const data = await projectApi.listProjects();
      setProjects(data);
    } catch (error) {
      message.error('Failed to fetch projects');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (activeTab === 'tenants') {
      fetchTenants();
    } else {
      fetchProjects();
    }
  }, [activeTab]);

  const handleCreate = () => {
    setEditingItem(null);
    form.resetFields();
    if (activeTab === 'tenants') {
      form.setFieldsValue({
        max_instances: 10,
        max_cpu: '10',
        max_memory: '20Gi',
        max_storage: '100Gi',
      });
    }
    setModalVisible(true);
  };

  const handleEdit = (record: Tenant | Project) => {
    setEditingItem(record);
    form.setFieldsValue(record);
    setModalVisible(true);
  };

  const handleDelete = async (id: string) => {
    if (activeTab === 'tenants') {
      try {
        await tenantApi.deleteTenant(id);
        message.success('Tenant deleted successfully');
        fetchTenants();
      } catch (error) {
        message.error('Failed to delete tenant');
      }
    } else {
      try {
        await projectApi.deleteProject(id);
        message.success('Project deleted successfully');
        fetchProjects();
      } catch (error) {
        message.error('Failed to delete project');
      }
    }
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      if (activeTab === 'tenants') {
        if (editingItem) {
          await tenantApi.updateTenant(editingItem.id, values);
          message.success('Tenant updated successfully');
        } else {
          await tenantApi.createTenant(values);
          message.success('Tenant created successfully');
        }
        fetchTenants();
      } else {
        if (editingItem) {
          await projectApi.updateProject(editingItem.id, values);
          message.success('Project updated successfully');
        } else {
          await projectApi.createProject(values);
          message.success('Project created successfully');
        }
        fetchProjects();
      }
      setModalVisible(false);
    } catch (error) {
      message.error('Failed to save item');
    }
  };

  const tenantColumns: ColumnsType<Tenant> = [
    {
      title: 'Name',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: 'Max Instances',
      dataIndex: 'max_instances',
      key: 'max_instances',
    },
    {
      title: 'Max CPU',
      dataIndex: 'max_cpu',
      key: 'max_cpu',
    },
    {
      title: 'Max Memory',
      dataIndex: 'max_memory',
      key: 'max_memory',
    },
    {
      title: 'Max Storage',
      dataIndex: 'max_storage',
      key: 'max_storage',
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
            title="Are you sure you want to delete this item?"
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

  const projectColumns: ColumnsType<Project> = [
    {
      title: 'Name',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: 'Tenant ID',
      dataIndex: 'tenant_id',
      key: 'tenant_id',
    },
    {
      title: 'Created At',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => new Date(date).toLocaleString(),
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
            title="Are you sure you want to delete this project?"
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
    <Card title="Tenant & Project Management">
      <Tabs
        activeKey={activeTab}
        onChange={setActiveTab}
        items={[
          {
            key: 'tenants',
            label: 'Tenants',
          },
          {
            key: 'projects',
            label: 'Projects',
          },
        ]}
      />
      <div style={{ marginBottom: 16 }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
          Create {activeTab === 'tenants' ? 'Tenant' : 'Project'}
        </Button>
      </div>
      {activeTab === 'tenants' ? (
        <Table
          columns={tenantColumns}
          dataSource={tenants}
          loading={loading}
          rowKey="id"
          pagination={{
            showSizeChanger: true,
            showQuickJumper: true,
          }}
        />
      ) : (
        <Table
          columns={projectColumns}
          dataSource={projects}
          loading={loading}
          rowKey="id"
          pagination={{
            showSizeChanger: true,
            showQuickJumper: true,
          }}
        />
      )}
      <Modal
        title={editingItem ? `Edit ${activeTab === 'tenants' ? 'Tenant' : 'Project'}` : `Create ${activeTab === 'tenants' ? 'Tenant' : 'Project'}`}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
        width={600}
      >
        <Form form={form} layout="vertical">
          {activeTab === 'tenants' ? (
            <>
              <Form.Item
                name="name"
                label="Name"
                rules={[{ required: true, message: 'Please input name' }]}
              >
                <Input placeholder="Tenant name" />
              </Form.Item>
              <Form.Item
                name="max_instances"
                label="Max Instances"
                rules={[{ required: true, message: 'Please input max instances' }]}
              >
                <InputNumber min={0} style={{ width: '100%' }} />
              </Form.Item>
              <Form.Item
                name="max_cpu"
                label="Max CPU"
              >
                <Input placeholder="10" />
              </Form.Item>
              <Form.Item
                name="max_memory"
                label="Max Memory"
              >
                <Input placeholder="20Gi" />
              </Form.Item>
              <Form.Item
                name="max_storage"
                label="Max Storage"
              >
                <Input placeholder="100Gi" />
              </Form.Item>
            </>
          ) : (
            <>
              <Form.Item
                name="tenant_id"
                label="Tenant ID"
                rules={[{ required: true, message: 'Please input tenant ID' }]}
              >
                <Input placeholder="Tenant ID" />
              </Form.Item>
              <Form.Item
                name="name"
                label="Name"
                rules={[{ required: true, message: 'Please input project name' }]}
              >
                <Input placeholder="Project name" />
              </Form.Item>
            </>
          )}
        </Form>
      </Modal>
    </Card>
  );
};

export default TenantProjectManagement;