import { Navigate } from 'react-router-dom';
import { useAuth } from '@/providers/AuthProvider';
import React from "react";

export function ProtectedRoute({ children }:{children: React.ReactNode}) {
    const { appUser, meInProgress } = useAuth();

    if (!appUser) {
        if (meInProgress) {
            return <div></div>;
        }
        return <Navigate to="/login" replace />;
    }

    return children;
}
