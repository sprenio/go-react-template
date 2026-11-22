import { Outlet, Navigate } from 'react-router-dom';
import { AuthenticatedLayout } from '@/layouts/AuthenticatedLayout';
import { useAuth } from '@/providers/AuthProvider';

export function ProtectedRoute() {
  const { appUser } = useAuth();

  if (!appUser) {
    return <Navigate to="/login" replace />;
  }

  return (
    <AuthenticatedLayout>
      <Outlet />
    </AuthenticatedLayout>
  );
}
