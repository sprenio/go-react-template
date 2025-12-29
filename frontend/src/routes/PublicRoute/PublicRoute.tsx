import { Navigate } from 'react-router-dom';
import { useAuth } from '@/providers/AuthProvider';
import React from "react";

export function PublicRoute({ children }: { children: React.ReactNode }) {
    const { token } = useAuth();
    if (token) return <Navigate to="/" replace />;
    return children;
}
