import { Navigate } from 'react-router-dom';
import { useAuth } from '@/providers/AuthProvider';
import React from "react";

export function PublicRoute({ children }: { children: React.ReactNode }) {
    const { appUser } = useAuth();
    if (appUser) return <Navigate to="/" replace />;
    return children;
}
