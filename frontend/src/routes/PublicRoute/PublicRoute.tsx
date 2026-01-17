import { PublicLayout } from '@/layouts/PublicLayout';
import { useAuth } from '@/providers/AuthProvider';
import { Navigate } from 'react-router-dom';

export function PublicRoute() {
  const { appUser } = useAuth();
  return appUser ?  <Navigate to="/dashboard" replace /> : <PublicLayout />;
}