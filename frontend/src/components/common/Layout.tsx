import React, { useState, useEffect } from 'react';
import { Layout, Menu, theme, Dropdown, Avatar } from 'antd';
import { AppstoreOutlined, SettingOutlined, ClusterOutlined, FileOutlined, UserOutlined, LogoutOutlined, SafetyOutlined } from '@ant-design/icons';
import { useNavigate, useLocation } from 'react-router-dom';
import { tokenManager, User } from '../../api/auth';

const { Header, Content } = Layout;

const { Sider } = Layout;

interface AppLayoutProps {
  children: React.ReactNode;
}

const AppLayout: React.FC<AppLayoutProps> = ({ children }) => {
  const [collapsed, setCollapsed] = useState(false);
  const [user, setUser] = useState<User | null>(null);
  const navigate = useNavigate();
  const location = useLocation();
  const {
    token: { colorBgContainer },
  } = theme.useToken();

  // Load user from localStorage on mount
  useEffect(() => {
    const currentUser = tokenManager.getUser();
    setUser(currentUser);
  }, []);

  // Handle logout
  const handleLogout = () => {
    tokenManager.clearTokens();
    navigate('/login');
  };

  // User menu items
  const userMenuItems = [
    {
      key: 'otp',
      icon: <SafetyOutlined />,
      label: '二次验证',
      onClick: () => navigate('/settings/otp'),
    },
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: '个人信息',
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
      onClick: handleLogout,
    },
  ];

  const menuItems = [
    {
      key: '/instances',
      icon: <ClusterOutlined />,
      label: '实例管理',
    },
    {
      key: '/configs',
      icon: <FileOutlined />,
      label: '配置管理',
    },
    {
      key: '/monitoring',
      icon: <AppstoreOutlined />,
      label: '监控中心',
    },
    {
      key: '/settings',
      icon: <SettingOutlined />,
      label: '系统设置',
    },
  ];

  const handleMenuClick = ({ key }: { key: string }) => {
    navigate(key);
  };

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider collapsible collapsed={collapsed} onCollapse={(value) => setCollapsed(value)}>
        <div style={{ height: 64, display: 'flex', alignItems: 'center', justifyContent: 'center', color: '#fff', fontSize: '18px', fontWeight: 'bold' }}>
          {collapsed ? 'OCC' : 'Open Cluster Claw'}
        </div>
        <Menu
          theme="dark"
          selectedKeys={[location.pathname]}
          mode="inline"
          items={menuItems}
          onClick={handleMenuClick}
        />
      </Sider>
      <Layout>
        <Header style={{ background: '#fff', padding: '0 24px', display: 'flex', justifyContent: 'space-between', alignItems: 'center', borderBottom: '1px solid #f0f0f0' }}>
          <div />
          <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
            <div style={{ display: 'flex', alignItems: 'center', cursor: 'pointer' }}>
              <Avatar icon={<UserOutlined />} style={{ marginRight: 8 }} />
              <span>{user?.username || 'Unknown'}</span>
              <span style={{ marginLeft: 12, color: '#999' }}>{user?.role || 'user'}</span>
            </div>
          </Dropdown>
        </Header>
        <Content style={{ padding: '24px' }}>
          <div
            style={{
              padding: 24,
              minHeight: 360,
              background: colorBgContainer,
              borderRadius: 8,
            }}
          >
            {children}
          </div>
        </Content>
      </Layout>
    </Layout>
  );
};

export default AppLayout;