import { AuthenticatedLayout } from '@/layouts/AuthenticatedLayout';
import { useAuth } from '@/providers/AuthProvider';
import { Navigate } from 'react-router-dom';

export function ProtectedRoute() {
  const {appUser} = useAuth();
  return appUser ? <AuthenticatedLayout /> : <Navigate to="/login" replace />;
}