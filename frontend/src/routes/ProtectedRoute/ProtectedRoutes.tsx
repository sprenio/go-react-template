import { Route, Navigate } from 'react-router-dom';
import { ProtectedRoute } from './ProtectedRoute';
import { Dashboard } from '@/pages/Dashboard';
import { Settings } from '@/pages/Settings';

export function getProtectedRoutes() {
  return (
    <Route element={<ProtectedRoute />}>
      <Route path="dashboard" element={<Dashboard />} />
      <Route path="settings" element={<Settings />} />
      <Route index element={<Navigate to="/dashboard" replace />} />
    </Route>
  );
}
