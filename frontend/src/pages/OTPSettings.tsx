import React, { useState, useEffect } from 'react';
import { Card, Steps, Button, Input, Modal, Alert, Space, Typography, List, Tag } from 'antd';
import { QrcodeOutlined, SafetyOutlined, CopyOutlined } from '@ant-design/icons';
import { otpApi } from '../api/otp';
import type { OTPStatus } from '../api/otp';

const { Title, Text, Paragraph } = Typography;

interface GenerateOTPResponse {
  secret: string;
  qr_code: string;
}

interface EnableOTPResponse {
  backup_codes: string[];
}

const OTPSettings: React.FC = () => {
  const [currentStep, setCurrentStep] = useState(0);
  const [loading, setLoading] = useState(false);
  const [otpStatus, setOtpStatus] = useState<OTPStatus>({ otp_enabled: false });

  // Generate OTP step
  const [secret, setSecret] = useState('');
  const [qrCode, setQrCode] = useState('');
  const [verifyCode, setVerifyCode] = useState('');

  // Backup codes
  const [backupCodes, setBackupCodes] = useState<string[]>([]);
  const [showBackupModal, setShowBackupModal] = useState(false);

  // Disable OTP step
  const [disableCode, setDisableCode] = useState('');

  // Fetch OTP status on mount
  useEffect(() => {
    fetchOTPStatus();
  }, []);

  const fetchOTPStatus = async () => {
    try {
      const status = await otpApi.getOTPStatus();
      setOtpStatus(status);
    } catch (error) {
      console.error('Failed to fetch OTP status:', error);
    }
  };

  const handleGenerateSecret = async () => {
    setLoading(true);
    try {
      const response: GenerateOTPResponse = await otpApi.generateSecret();
      setSecret(response.secret);
      setQrCode(response.qr_code);
      setCurrentStep(1);
    } catch (error: any) {
      Modal.error({
        title: '生成失败',
        content: error.response?.data?.message || '生成 OTP 密钥失败',
      });
    } finally {
      setLoading(false);
    }
  };

  const handleEnableOTP = async () => {
    if (!verifyCode || verifyCode.length !== 6) {
      Modal.error({
        title: '验证码格式错误',
        content: '请输入 6 位验证码',
      });
      return;
    }

    setLoading(true);
    try {
      const response: EnableOTPResponse = await otpApi.enableOTP({ code: verifyCode });
      setBackupCodes(response.backup_codes);
      setCurrentStep(2);
      await fetchOTPStatus();
    } catch (error: any) {
      Modal.error({
        title: '启用失败',
        content: error.response?.data?.message || '启用 OTP 失败',
      });
    } finally {
      setLoading(false);
    }
  };

  const handleDisableOTP = async () => {
    if (!disableCode || disableCode.length !== 6) {
      Modal.error({
        title: '验证码格式错误',
        content: '请输入 6 位验证码',
      });
      return;
    }

    setLoading(true);
    try {
      await otpApi.disableOTP({ code: disableCode });
      setCurrentStep(0);
      setSecret('');
      setQrCode('');
      setVerifyCode('');
      setDisableCode('');
      await fetchOTPStatus();
      Modal.success({
        title: '已禁用',
        content: 'OTP 二次验证已成功禁用',
      });
    } catch (error: any) {
      Modal.error({
        title: '禁用失败',
        content: error.response?.data?.message || '禁用 OTP 失败',
      });
    } finally {
      setLoading(false);
    }
  };

  const handleShowBackupCodes = async () => {
    try {
      const codes = await otpApi.getBackupCodes();
      setBackupCodes(codes);
      setShowBackupModal(true);
    } catch (error: any) {
      Modal.error({
        title: '获取失败',
        content: error.response?.data?.message || '获取备份码失败',
      });
    }
  };

  const handleCopyCode = (code: string) => {
    navigator.clipboard.writeText(code);
    Modal.success({
      title: '已复制',
      content: `备份码 ${code} 已复制到剪贴板`,
    });
  };

  const stepsItems = [
    {
      title: '选择操作',
    },
    {
      title: '扫描二维码',
    },
    {
      title: '启用成功',
    },
    {
      title: '禁用 OTP',
    },
  ];

  // Render content based on current step
  const renderContent = () => {
    switch (currentStep) {
      case 0:
        return (
          <div style={{ padding: '20px 0' }}>
            <Space direction="vertical" size="large" style={{ width: '100%' }}>
              {!otpStatus.otp_enabled ? (
                <Button
                  type="primary"
                  icon={<QrcodeOutlined />}
                  onClick={handleGenerateSecret}
                  loading={loading}
                  size="large"
                  block
                >
                  启用 OTP 验证
                </Button>
              ) : (
                <Space direction="vertical" size="middle" style={{ width: '100%' }}>
                  <Button
                    icon={<SafetyOutlined />}
                    onClick={() => {
                      setDisableCode('');
                      setCurrentStep(3);
                    }}
                    size="large"
                    block
                  >
                    禁用 OTP 验证
                  </Button>
                  <Button
                    onClick={handleShowBackupCodes}
                    size="large"
                    block
                  >
                    查看备份码
                  </Button>
                </Space>
              )}
            </Space>
          </div>
        );

      case 1:
        return (
          <div style={{ padding: '20px 0', textAlign: 'center' }}>
            <Space direction="vertical" size="large" style={{ width: '100%' }}>
              <Alert
                message="使用 Google Authenticator 等应用扫描下方二维码"
                type="info"
                showIcon
              />
              <div style={{ display: 'flex', justifyContent: 'center', padding: '20px' }}>
                <img src={qrCode} alt="QR Code" style={{ width: 256, height: 256 }} />
              </div>
              <div>
                <Text strong>密钥：</Text>
                <Input
                  value={secret}
                  readOnly
                  suffix={
                    <CopyOutlined
                      style={{ cursor: 'pointer' }}
                      onClick={() => {
                        navigator.clipboard.writeText(secret);
                        Modal.success({ title: '已复制', content: '密钥已复制到剪贴板' });
                      }}
                    />
                  }
                  style={{ fontFamily: 'monospace', fontSize: '16px' }}
                />
              </div>
              <div>
                <Text>1. 使用 Google Authenticator、Authy 等 OTP 应用</Text>
                <Text>2. 扫描上方二维码或手动输入密钥</Text>
                <Text>3. 输入应用显示的 6 位验证码</Text>
              </div>
              <Input
                placeholder="请输入应用显示的 6 位验证码"
                maxLength={6}
                value={verifyCode}
                onChange={(e) => setVerifyCode(e.target.value)}
                size="large"
                style={{ textAlign: 'center', letterSpacing: '8px', fontSize: '24px' }}
              />
              <Space style={{ width: '100%' }}>
                <Button onClick={() => setCurrentStep(0)} style={{ flex: 1 }}>
                  取消
                </Button>
                <Button
                  type="primary"
                  onClick={handleEnableOTP}
                  loading={loading}
                  disabled={!verifyCode || verifyCode.length !== 6}
                  style={{ flex: 1 }}
                >
                  验证并启用
                </Button>
              </Space>
            </Space>
          </div>
        );

      case 2:
        return (
          <div style={{ padding: '20px 0' }}>
            <Space direction="vertical" size="large" style={{ width: '100%' }}>
              <Alert
                message="OTP 二次验证已成功启用！"
                description="请保存以下备份码，当无法使用 Authenticator 时可以使用这些代码登录。"
                type="success"
                showIcon
              />
              <List
                header={
                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <Text strong>备份码（每个仅可使用一次）</Text>
                    <Button type="link" onClick={() => setShowBackupModal(true)}>
                      下载
                    </Button>
                  </div>
                }
                dataSource={backupCodes}
                renderItem={(code, index) => (
                  <List.Item
                    actions={[
                      <Button
                        icon={<CopyOutlined />}
                        type="text"
                        onClick={() => handleCopyCode(code)}
                      >
                        复制
                      </Button>,
                    ]}
                  >
                    <Space>
                      <Tag color="blue">{index + 1}</Tag>
                      <Text style={{ fontFamily: 'monospace', fontSize: '16px', letterSpacing: '2px' }}>
                        {code}
                      </Text>
                    </Space>
                  </List.Item>
                )}
              />
              <Button type="primary" onClick={() => setCurrentStep(0)} block>
                完成
              </Button>
            </Space>
          </div>
        );

      case 3:
        return (
          <div style={{ padding: '20px 0' }}>
            <Space direction="vertical" size="large" style={{ width: '100%' }}>
              <Alert
                message="确定要禁用 OTP 二次验证吗？"
                description="禁用后，账户安全性将降低。"
                type="warning"
                showIcon
              />
              <Paragraph>
                请输入 Authenticator 应用显示的 6 位验证码以确认禁用：
              </Paragraph>
              <Input
                placeholder="请输入 6 位验证码"
                maxLength={6}
                value={disableCode}
                onChange={(e) => setDisableCode(e.target.value)}
                size="large"
                style={{ textAlign: 'center', letterSpacing: '8px', fontSize: '24px' }}
              />
              <Space style={{ width: '100%' }}>
                <Button onClick={() => setCurrentStep(0)} style={{ flex: 1 }}>
                  取消
                </Button>
                <Button
                  type="primary"
                  danger
                  onClick={handleDisableOTP}
                  loading={loading}
                  disabled={!disableCode || disableCode.length !== 6}
                  style={{ flex: 1 }}
                >
                  确认禁用
                </Button>
              </Space>
            </Space>
          </div>
        );

      default:
        return null;
    }
  };

  return (
    <div style={{ padding: '24px' }}>
      <Card title="二次验证 (2FA)">
        <div style={{ marginBottom: '24px' }}>
          <Space>
            <SafetyOutlined style={{ fontSize: '24px' }} />
            <div>
              <Title level={5} style={{ margin: 0 }}>
                OTP 二次验证
              </Title>
              <Text type={otpStatus.otp_enabled ? 'success' : 'secondary'}>
                {otpStatus.otp_enabled ? '已启用' : '未启用'}
              </Text>
            </div>
          </Space>
        </div>

        <Steps current={currentStep} items={stepsItems} />

        <div style={{ marginTop: '24px' }}>
          {renderContent()}
        </div>
      </Card>

      <Modal
        title="备份码"
        open={showBackupModal}
        onCancel={() => setShowBackupModal(false)}
        footer={[
          <Button key="close" onClick={() => setShowBackupModal(false)}>
            关闭
          </Button>,
          <Button
            key="download"
            type="primary"
            onClick={() => {
              const text = backupCodes.join('\n');
              const blob = new Blob([text], { type: 'text/plain' });
              const url = URL.createObjectURL(blob);
              const a = document.createElement('a');
              a.href = url;
              a.download = 'otp-backup-codes.txt';
              a.click();
              URL.revokeObjectURL(url);
            }}
          >
            下载
          </Button>,
        ]}
        width={600}
      >
        <Alert
          message="重要提示"
          description="请妥善保管这些备份码，每个备份码只能使用一次。建议将它们保存在安全的地方。"
          type="warning"
          showIcon
          style={{ marginBottom: '16px' }}
        />
        <List
          dataSource={backupCodes}
          renderItem={(code, index) => (
            <List.Item
              actions={[
                <Button
                  icon={<CopyOutlined />}
                  type="text"
                  onClick={() => handleCopyCode(code)}
                >
                  复制
                </Button>,
              ]}
            >
              <Space>
                <Tag color="blue">{index + 1}</Tag>
                <Text style={{ fontFamily: 'monospace', fontSize: '16px', letterSpacing: '2px' }}>
                  {code}
                </Text>
              </Space>
            </List.Item>
          )}
        />
      </Modal>
    </div>
  );
};

export default OTPSettings;
