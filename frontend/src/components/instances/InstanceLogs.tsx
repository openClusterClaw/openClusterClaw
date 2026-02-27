import React, { useState, useEffect } from 'react';
import { Card, Button, Space, InputNumber, message, Empty, Typography, Spin } from 'antd';
import { ReloadOutlined, CopyOutlined } from '@ant-design/icons';
import { instanceApi } from '@/api/instance';

const { Text } = Typography;

interface InstanceLogsProps {
  instanceId: string;
}

const InstanceLogs: React.FC<InstanceLogsProps> = ({ instanceId }) => {
  const [logs, setLogs] = useState('');
  const [tailLines, setTailLines] = useState(100);
  const [loading, setLoading] = useState(false);
  const [autoRefresh, setAutoRefresh] = useState(false);

  const fetchLogs = async () => {
    setLoading(true);
    try {
      const data = await instanceApi.getInstanceLogs(instanceId, tailLines);
      setLogs(data);
    } catch (error: any) {
      if (error.response?.status === 500) {
        setLogs('Logs not available: ' + (error.response.data?.message || 'K8S integration not enabled'));
      } else {
        message.error('Failed to fetch logs');
      }
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchLogs();
  }, [instanceId]);

  useEffect(() => {
    if (autoRefresh) {
      const interval = setInterval(fetchLogs, 5000);
      return () => clearInterval(interval);
    }
  }, [autoRefresh, tailLines]);

  const handleRefresh = () => {
    fetchLogs();
  };

  const handleCopy = () => {
    navigator.clipboard.writeText(logs);
    message.success('Logs copied to clipboard');
  };

  return (
    <Card
      title="Instance Logs"
      extra={
        <Space>
          <InputNumber
            min={10}
            max={1000}
            value={tailLines}
            onChange={(value) => setTailLines(value || 100)}
            style={{ width: 120 }}
            addonAfter="lines"
          />
          <Button
            type={autoRefresh ? 'primary' : 'default'}
            onClick={() => setAutoRefresh(!autoRefresh)}
          >
            {autoRefresh ? 'Auto Refresh On' : 'Auto Refresh Off'}
          </Button>
          <Button icon={<ReloadOutlined />} onClick={handleRefresh}>
            Refresh
          </Button>
          {logs && <Button icon={<CopyOutlined />} onClick={handleCopy}>Copy</Button>}
        </Space>
      }
    >
      <Spin spinning={loading}>
        {logs ? (
          <pre
            style={{
              background: '#1a1a1a',
              color: '#fff',
              padding: '16px',
              borderRadius: '4px',
              maxHeight: '500px',
              overflow: 'auto',
              fontFamily: 'monospace',
              fontSize: '12px',
              whiteSpace: 'pre-wrap',
              wordBreak: 'break-all',
            }}
          >
            {logs || <Text type="secondary">No logs available</Text>}
          </pre>
        ) : (
          <Empty description="No logs available" />
        )}
      </Spin>
    </Card>
  );
};

export default InstanceLogs;