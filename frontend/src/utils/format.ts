import type { InstanceStatus } from '@/types';

export const getStatusColor = (status: InstanceStatus): string => {
  const colorMap: Record<InstanceStatus, string> = {
    Creating: 'blue',
    Running: 'green',
    Stopped: 'default',
    Failed: 'red',
    Destroyed: 'default',
  };
  return colorMap[status] || 'default';
};

export const formatDateTime = (dateString: string): string => {
  const date = new Date(dateString);
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  });
};

export const formatResourceSize = (value: string): string => {
  // Format resource size (e.g., "4Gi" -> "4 GB")
  if (value.endsWith('Gi')) {
    return value.replace('Gi', ' GB');
  }
  if (value.endsWith('Mi')) {
    return value.replace('Mi', ' MB');
  }
  return value;
};