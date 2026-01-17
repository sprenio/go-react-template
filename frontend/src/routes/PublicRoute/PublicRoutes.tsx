import { PublicRoute } from './PublicRoute';
import { Navigate, Route } from 'react-router-dom';
import { Login } from '@/pages/Login';
import { Register } from '@/pages/Register';
import { ChangePassword } from '@/pages/ChangePassword';
import { ResetPassword } from '@/pages/ResetPassword';
import { Confirm } from '@/pages/Confirm';

export function getPublicRoutes() {
  return (
    <Route element={<PublicRoute />}>
      <Route path="login" element={<Login />} />
      <Route path="register" element={<Register />} />
      <Route path="reset-password/:hash" element={<ChangePassword />} />
      <Route path="reset-password" element={<ResetPassword />} />
      <Route path="confirm/:hash" element={<Confirm />} />
      <Route index element={<Navigate to="/login" replace />} />
    </Route>
  );
}