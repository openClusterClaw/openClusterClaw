import React, { useState } from 'react';
import { Form, Input, Button, Card, message, Steps, Alert } from 'antd';
import { UserOutlined, LockOutlined, SafetyOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { authApi, tokenManager } from '../api/auth';
import { otpApi } from '../api/otp';
import type { LoginOTPResponse } from '../api/auth';

const { Step } = Steps;

const Login: React.FC = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [currentStep, setCurrentStep] = useState(0);
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [otpCode, setOtpCode] = useState('');
  const [tempToken, setTempToken] = useState<string | null>(null);

  const handleLogin = async (values: { username: string; password: string }) => {
    setLoading(true);
    try {
      const response: LoginOTPResponse = await authApi.login(values);

      // Check if response is valid
      if (!response) {
        throw new Error('Invalid response from server');
      }

      // Check if OTP is required
      if (response.requires_otp) {
        if (response.temp_token) {
          setTempToken(response.temp_token);
          setCurrentStep(1);
          message.info('请输入 OTP 验证码');
        } else {
          throw new Error('OTP required but no temp token provided');
        }
      } else if (response.access_token && response.refresh_token && response.user) {
        // OTP not required, complete login
        tokenManager.setTokens(
          response.access_token,
          response.refresh_token,
          response.user
        );
        message.success('Login successful');
        navigate('/');
      } else {
        throw new Error('Invalid response format');
      }
    } catch (error: any) {
      console.error('Login error:', error);
      message.error(error.response?.data?.message || error.message || 'Login failed');
    } finally {
      setLoading(false);
    }
  };

  const handleVerifyOTP = async () => {
    if (!otpCode || otpCode.length !== 6) {
      message.error('请输入 6 位验证码');
      return;
    }

    if (!tempToken) {
      message.error('会话已过期，请重新登录');
      setCurrentStep(0);
      return;
    }

    setLoading(true);
    try {
      const response = await otpApi.verifyOTP({
        temp_token: tempToken,
        code: otpCode,
      });

      // Check if response is valid
      if (!response || !response.access_token || !response.user) {
        throw new Error('Invalid response from server');
      }

      // Complete login with verified OTP
      tokenManager.setTokens(
        response.access_token,
        response.refresh_token,
        response.user
      );
      message.success('Login successful');
      navigate('/');
    } catch (error: any) {
      console.error('OTP verification error:', error);
      message.error(error.response?.data?.message || error.message || 'OTP verification failed');
    } finally {
      setLoading(false);
    }
  };

  const handleBack = () => {
    setCurrentStep(0);
    setTempToken(null);
    setOtpCode('');
  };

  return (
    <div style={{
      display: 'flex',
      justifyContent: 'center',
      alignItems: 'center',
      minHeight: '100vh',
      background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    }}>
      <Card
        title="Open Cluster Claw"
        style={{
          width: currentStep === 0 ? 400 : 500,
          boxShadow: '0 8px 32px rgba(0, 0, 0, 0.1)',
        }}
      >
        <Steps current={currentStep} style={{ marginBottom: '24px' }} size="small">
          <Step title="登录" icon={<UserOutlined />} />
          <Step title="OTP 验证" icon={<SafetyOutlined />} />
        </Steps>

        {currentStep === 0 && (
          <Form
            name="login"
            onFinish={handleLogin}
            autoComplete="off"
            size="large"
          >
            <Form.Item
              name="username"
              rules={[{ required: true, message: 'Please input your username!' }]}
            >
              <Input
                prefix={<UserOutlined />}
                placeholder="Username"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
              />
            </Form.Item>

            <Form.Item
              name="password"
              rules={[{ required: true, message: 'Please input your password!' }]}
            >
              <Input.Password
                prefix={<LockOutlined />}
                placeholder="Password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
            </Form.Item>

            <Form.Item>
              <Button
                type="primary"
                htmlType="submit"
                loading={loading}
                block
              >
                Login
              </Button>
            </Form.Item>

            <div style={{ textAlign: 'center', marginTop: '16px', color: '#666' }}>
              <p>默认账号：admin / admin123</p>
            </div>
          </Form>
        )}

        {currentStep === 1 && (
          <div>
            <Alert
              message="OTP 二次验证"
              description="请输入 Google Authenticator 等应用显示的 6 位验证码"
              type="info"
              showIcon
              style={{ marginBottom: '24px' }}
            />
            <Form size="large">
              <Form.Item
                rules={[
                  { required: true, message: '请输入验证码' },
                  { len: 6, message: '验证码必须为 6 位' },
                ]}
              >
                <Input
                  prefix={<SafetyOutlined />}
                  placeholder="请输入 6 位验证码"
                  maxLength={6}
                  value={otpCode}
                  onChange={(e) => setOtpCode(e.target.value)}
                  style={{ textAlign: 'center', letterSpacing: '8px', fontSize: '24px' }}
                />
              </Form.Item>
            </Form>

            <div style={{ marginTop: '16px' }}>
              <Button
                type="primary"
                onClick={handleVerifyOTP}
                loading={loading}
                block
                style={{ marginBottom: '8px' }}
              >
                验证并登录
              </Button>
              <Button
                onClick={handleBack}
                block
              >
                返回登录
              </Button>
            </div>
          </div>
        )}
      </Card>
    </div>
  );
};

export default Login;
