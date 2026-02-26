import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { Card, Descriptions, Button, Space, Tag, message, Spin } from 'antd';
import { ArrowLeftOutlined, PlayCircleOutlined, StopOutlined, SyncOutlined, DeleteOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { instanceApi } from '@/api/instance';
import type { ClawInstance } from '@/types';
import { getStatusColor, formatDateTime, formatResourceSize } from '@/utils/format';

const InstanceDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [instance, setInstance] = useState<ClawInstance | null>(null);
  const [loading, setLoading] = useState(false);
  const [actionLoading, setActionLoading] = useState(false);
  const navigate = useNavigate();

  const fetchInstance = async () => {
    if (!id) return;
    setLoading(true);
    try {
      const { data } = await instanceApi.get(id);
      if (data.code === 0 && data.data) {
        setInstance(data.data);
      }
    } catch (error) {
      message.error('获取实例详情失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchInstance();
  }, [id]);

  const handleStart = async () => {
    if (!id) return;
    setActionLoading(true);
    try {
      const { data } = await instanceApi.start(id);
      if (data.code === 0) {
        message.success('实例启动成功');
        fetchInstance();
      }
    } catch (error) {
      message.error('实例启动失败');
    } finally {
      setActionLoading(false);
    }
  };

  const handleStop = async () => {
    if (!id) return;
    setActionLoading(true);
    try {
      const { data } = await instanceApi.stop(id);
      if (data.code === 0) {
        message.success('实例停止成功');
        fetchInstance();
      }
    } catch (error) {
      message.error('实例停止失败');
    } finally {
      setActionLoading(false);
    }
  };

  const handleRestart = async () => {
    if (!id) return;
    setActionLoading(true);
    try {
      const { data } = await instanceApi.restart(id);
      if (data.code === 0) {
        message.success('实例重启成功');
        fetchInstance();
      }
    } catch (error) {
      message.error('实例重启失败');
    } finally {
      setActionLoading(false);
    }
  };

  if (loading || !instance) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: 400 }}>
        <Spin size="large" />
      </div>
    );
  }

  return (
    <div>
      <Space style={{ marginBottom: 16 }}>
        <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/instances')}>
          返回
        </Button>
        <Space>
          {instance.status === 'Stopped' || instance.status === 'Failed' ? (
            <Button type="primary" icon={<PlayCircleOutlined />} loading={actionLoading} onClick={handleStart}>
              启动
            </Button>
          ) : instance.status === 'Running' ? (
            <Button icon={<StopOutlined />} loading={actionLoading} onClick={handleStop}>
              停止
            </Button>
          ) : null}
          <Button icon={<SyncOutlined />} loading={actionLoading} onClick={handleRestart}>
            重启
          </Button>
          <Button danger icon={<DeleteOutlined />}>
            删除
          </Button>
        </Space>
      </Space>

      <Card title="实例详情" bordered={false}>
        <Descriptions column={2} bordered>
          <Descriptions.Item label="实例名称">{instance.name}</Descriptions.Item>
          <Descriptions.Item label="实例ID">{instance.id}</Descriptions.Item>
          <Descriptions.Item label="类型">{instance.type}</Descriptions.Item>
          <Descriptions.Item label="版本">{instance.version}</Descriptions.Item>
          <Descriptions.Item label="状态">
            <Tag color={getStatusColor(instance.status)}>{instance.status}</Tag>
          </Descriptions.Item>
          <Descriptions.Item label="租户ID">{instance.tenant_id}</Descriptions.Item>
          <Descriptions.Item label="项目ID">{instance.project_id}</Descriptions.Item>
          <Descriptions.Item label="创建时间">{formatDateTime(instance.created_at)}</Descriptions.Item>
          <Descriptions.Item label="更新时间">{formatDateTime(instance.updated_at)}</Descriptions.Item>
          <Descriptions.Item label="CPU">{formatResourceSize(instance.resources?.cpu || '-')}</Descriptions.Item>
          <Descriptions.Item label="内存">{formatResourceSize(instance.resources?.memory || '-')}</Descriptions.Item>
        </Descriptions>

        {instance.storage && (
          <Card title="存储配置" style={{ marginTop: 16 }}>
            <Descriptions column={1} bordered>
              <Descriptions.Item label="配置目录">{instance.storage.config_dir}</Descriptions.Item>
              <Descriptions.Item label="数据目录">{instance.storage.data_dir}</Descriptions.Item>
              <Descriptions.Item label="存储大小">{formatResourceSize(instance.storage.size)}</Descriptions.Item>
            </Descriptions>
          </Card>
        )}
      </Card>
    </div>
  );
};

export default InstanceDetail;