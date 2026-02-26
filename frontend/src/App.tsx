import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { Layout } from 'antd';
import AppLayout from './components/common/Layout';
import InstanceList from './pages/InstanceList';
import InstanceDetail from './pages/InstanceDetail';

const App: React.FC = () => {
  return (
    <Router>
      <Layout style={{ minHeight: '100vh' }}>
        <AppLayout>
          <Routes>
            <Route path="/" element={<Navigate to="/instances" replace />} />
            <Route path="/instances" element={<InstanceList />} />
            <Route path="/instances/:id" element={<InstanceDetail />} />
          </Routes>
        </AppLayout>
      </Layout>
    </Router>
  );
};

export default App;