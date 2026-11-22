import { Route, Navigate } from 'react-router-dom';
import { Dashboard } from '@/pages/Dashboard';
import { Settings } from '@/pages/Settings';
import { NotFound } from '@/pages/NotFound';
import { Confirm } from '@/pages/Confirm';
import { ProtectedRoute } from './ProtectedRoute';

export function getProtectedRoutes() {
  return (
    <Route element={<ProtectedRoute />}>
      <Route path="dashboard" element={<Dashboard />} />
      <Route path="settings" element={<Settings />} />
      <Route path="confirm/:hash" element={<Confirm />} />
      <Route index element={<Navigate to="/dashboard" replace />} />
      <Route path="*" element={<NotFound />} />
    </Route>
  );
}
