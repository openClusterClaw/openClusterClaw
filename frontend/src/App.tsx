import { BrowserRouter as Router, Routes, Route, Navigate, useLocation } from 'react-router-dom';
import { Layout } from 'antd';
import AppLayout from './components/common/Layout';
import InstanceList from './pages/InstanceList';
import InstanceDetail from './pages/InstanceDetail';
import Login from './pages/Login';
import { tokenManager } from './api/auth';

// ProtectedRoute component that checks authentication
const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const location = useLocation();

  if (!tokenManager.isAuthenticated()) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  return <>{children}</>;
};

const App: React.FC = () => {
  return (
    <Router>
      <Routes>
        {/* Public routes */}
        <Route path="/login" element={<Login />} />

        {/* Protected routes */}
        <Route
          path="/*"
          element={
            <ProtectedRoute>
              <Layout style={{ minHeight: '100vh' }}>
                <AppLayout>
                  <Routes>
                    <Route path="/" element={<Navigate to="/instances" replace />} />
                    <Route path="/instances" element={<InstanceList />} />
                    <Route path="/instances/:id" element={<InstanceDetail />} />
                  </Routes>
                </AppLayout>
              </Layout>
            </ProtectedRoute>
          }
        />
      </Routes>
    </Router>
  );
};

export default App;
