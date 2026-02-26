import React from 'react';
import { Modal, Form, Input, Select, InputNumber, message } from 'antd';
import { instanceApi } from '@/api/instance';
import type { CreateInstanceRequest } from '@/types';

interface CreateInstanceModalProps {
  open: boolean;
  onClose: () => void;
  onSuccess: () => void;
}

const CreateInstanceModal: React.FC<CreateInstanceModalProps> = ({ open, onClose, onSuccess }) => {
  const [form] = Form.useForm();
  const [loading, setLoading] = React.useState(false);

  const handleOk = async () => {
    try {
      const values = await form.validateFields();
      setLoading(true);

      const data: CreateInstanceRequest = {
        name: values.name,
        tenant_id: values.tenant_id,
        project_id: values.project_id,
        type: values.type,
        version: values.version,
        cpu: values.cpu ? `${values.cpu}m` : undefined,
        memory: values.memory ? `${values.memory}Mi` : undefined,
      };

      const response = await instanceApi.create(data);

      if (response.data.code === 0) {
        message.success('实例创建成功');
        form.resetFields();
        onSuccess();
        onClose();
      } else {
        message.error(response.data.message || '创建失败');
      }
    } catch (error) {
      message.error('创建实例失败');
    } finally {
      setLoading(false);
    }
  };

  const handleCancel = () => {
    form.resetFields();
    onClose();
  };

  return (
    <Modal
      title="创建实例"
      open={open}
      onOk={handleOk}
      onCancel={handleCancel}
      confirmLoading={loading}
      width={600}
    >
      <Form form={form} layout="vertical">
        <Form.Item
          name="name"
          label="实例名称"
          rules={[{ required: true, message: '请输入实例名称' }]}
        >
          <Input placeholder="请输入实例名称" />
        </Form.Item>

        <Form.Item
          name="tenant_id"
          label="租户ID"
          rules={[{ required: true, message: '请输入租户ID' }]}
        >
          <Input placeholder="请输入租户ID" />
        </Form.Item>

        <Form.Item
          name="project_id"
          label="项目ID"
          rules={[{ required: true, message: '请输入项目ID' }]}
        >
          <Input placeholder="请输入项目ID" />
        </Form.Item>

        <Form.Item
          name="type"
          label="实例类型"
          rules={[{ required: true, message: '请选择实例类型' }]}
        >
          <Select placeholder="请选择实例类型">
            <Select.Option value="OpenClaw">OpenClaw</Select.Option>
            <Select.Option value="NanoClaw">NanoClaw</Select.Option>
          </Select>
        </Form.Item>

        <Form.Item
          name="version"
          label="版本"
          rules={[{ required: true, message: '请输入版本' }]}
        >
          <Input placeholder="请输入版本号，如 1.0.0" />
        </Form.Item>

        <Form.Item name="cpu" label="CPU (m)">
          <InputNumber min={100} max={8000} step={100} style={{ width: '100%' }} placeholder="100-8000" />
        </Form.Item>

        <Form.Item name="memory" label="内存 (Mi)">
          <InputNumber min={128} max={16384} step={128} style={{ width: '100%' }} placeholder="128-16384" />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default CreateInstanceModal;