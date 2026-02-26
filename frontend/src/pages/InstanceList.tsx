import React, { useEffect, useState } from 'react';
import { Table, Button, Space, Tag, message, Popconfirm, Input, Select } from 'antd';
import { PlusOutlined, ReloadOutlined, PlayCircleOutlined, StopOutlined, DeleteOutlined, SyncOutlined } from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import { useNavigate } from 'react-router-dom';
import { instanceApi } from '@/api/instance';
import type { ClawInstance } from '@/types';
import { getStatusColor, formatDateTime } from '@/utils/format';
import CreateInstanceModal from '@/components/instances/CreateInstanceModal';

const InstanceList: React.FC = () => {
  const [instances, setInstances] = useState<ClawInstance[]>([]);
  const [loading, setLoading] = useState(false);
  const [createModalOpen, setCreateModalOpen] = useState(false);
  const [actionLoading, setActionLoading] = useState<Record<string, boolean>>({});
  const navigate = useNavigate();

  const fetchInstances = async () => {
    setLoading(true);
    try {
      const { data } = await instanceApi.list();
      if (data.code === 0 && data.data) {
        setInstances(data.data.instances);
      }
    } catch (error) {
      message.error('获取实例列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchInstances();
  }, []);

  const handleStart = async (id: string) => {
    setActionLoading((prev) => ({ ...prev, [id]: true }));
    try {
      const { data } = await instanceApi.start(id);
      if (data.code === 0) {
        message.success('实例启动成功');
        fetchInstances();
      } else {
        message.error(data.message || '启动失败');
      }
    } catch (error) {
      message.error('实例启动失败');
    } finally {
      setActionLoading((prev) => ({ ...prev, [id]: false }));
    }
  };

  const handleStop = async (id: string) => {
    setActionLoading((prev) => ({ ...prev, [id]: true }));
    try {
      const { data } = await instanceApi.stop(id);
      if (data.code === 0) {
        message.success('实例停止成功');
        fetchInstances();
      } else {
        message.error(data.message || '停止失败');
      }
    } catch (error) {
      message.error('实例停止失败');
    } finally {
      setActionLoading((prev) => ({ ...prev, [id]: false }));
    }
  };

  const handleRestart = async (id: string) => {
    setActionLoading((prev) => ({ ...prev, [id]: true }));
    try {
      const { data } = await instanceApi.restart(id);
      if (data.code === 0) {
        message.success('实例重启成功');
        fetchInstances();
      } else {
        message.error(data.message || '重启失败');
      }
    } catch (error) {
      message.error('实例重启失败');
    } finally {
      setActionLoading((prev) => ({ ...prev, [id]: false }));
    }
  };

  const handleDelete = async (id: string) => {
    setActionLoading((prev) => ({ ...prev, [id]: true }));
    try {
      const { data } = await instanceApi.delete(id);
      if (data.code === 0) {
        message.success('实例删除成功');
        fetchInstances();
      } else {
        message.error(data.message || '删除失败');
      }
    } catch (error) {
      message.error('实例删除失败');
    } finally {
      setActionLoading((prev) => ({ ...prev, [id]: false }));
    }
  };

  const columns: ColumnsType<ClawInstance> = [
    {
      title: '实例名称',
      dataIndex: 'name',
      key: 'name',
      render: (text, record) => (
        <a onClick={() => navigate(`/instances/${record.id}`)}>{text}</a>
      ),
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
    },
    {
      title: '版本',
      dataIndex: 'version',
      key: 'version',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status) => <Tag color={getStatusColor(status)}>{status}</Tag>,
    },
    {
      title: 'CPU',
      dataIndex: ['resources', 'cpu'],
      key: 'cpu',
    },
    {
      title: '内存',
      dataIndex: ['resources', 'memory'],
      key: 'memory',
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (time) => formatDateTime(time),
    },
    {
      title: '操作',
      key: 'action',
      render: (_, record) => (
        <Space size="small">
          {record.status === 'Stopped' || record.status === 'Failed' ? (
            <Button
              type="link"
              size="small"
              icon={<PlayCircleOutlined />}
              loading={actionLoading[record.id]}
              onClick={() => handleStart(record.id)}
            >
              启动
            </Button>
          ) : record.status === 'Running' ? (
            <Button
              type="link"
              size="small"
              icon={<StopOutlined />}
              loading={actionLoading[record.id]}
              onClick={() => handleStop(record.id)}
            >
              停止
            </Button>
          ) : null}
          <Button
            type="link"
            size="small"
            icon={<SyncOutlined />}
            loading={actionLoading[record.id]}
            onClick={() => handleRestart(record.id)}
          >
            重启
          </Button>
          <Popconfirm
            title="确认删除"
            description="删除后实例将无法恢复"
            onConfirm={() => handleDelete(record.id)}
            okText="确认"
            cancelText="取消"
          >
            <Button
              type="link"
              size="small"
              danger
              icon={<DeleteOutlined />}
              loading={actionLoading[record.id]}
            >
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Space>
          <Input.Search placeholder="搜索实例" style={{ width: 200 }} />
          <Select placeholder="选择租户" style={{ width: 150 }} allowClear />
          <Select placeholder="选择状态" style={{ width: 150 }} allowClear />
        </Space>
        <Space>
          <Button icon={<ReloadOutlined />} onClick={fetchInstances}>
            刷新
          </Button>
          <Button type="primary" icon={<PlusOutlined />} onClick={() => setCreateModalOpen(true)}>
            创建实例
          </Button>
        </Space>
      </div>

      <Table
        columns={columns}
        dataSource={instances}
        rowKey="id"
        loading={loading}
        pagination={{
          showSizeChanger: true,
          showTotal: (total) => `共 ${total} 条`,
        }}
      />

      <CreateInstanceModal open={createModalOpen} onClose={() => setCreateModalOpen(false)} onSuccess={fetchInstances} />
    </div>
  );
};

export default InstanceList;