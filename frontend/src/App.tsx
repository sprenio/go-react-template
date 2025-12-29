import { Routes, Route, Navigate } from 'react-router-dom';
import {PublicLayout} from '@/layouts/PublicLayout';
import {AuthenticatedLayout} from '@/layouts/AuthenticatedLayout';
import {PublicRoute} from '@/routes/PublicRoute';
import {ProtectedRoute} from '@/routes/ProtectedRoute';
import {Login} from '@/pages/Login';
import {Register} from '@/pages/Register';
import {ChangePassword} from '@/pages/ChangePassword';
import {ResetPassword} from '@/pages/ResetPassword';
import {Confirm} from '@/pages/Confirm';
import {Dashboard} from '@/pages/Dashboard';
import {Settings} from '@/pages/Settings';
import {NotFound} from '@/pages/NotFound';

function App() {
    return (
        <Routes>
            {/* Public routes (dostępne tylko jeśli NIE jesteś zalogowany) */}
            <Route element={<PublicLayout />}>
                <Route
                    path="/login"
                    element={
                        <PublicRoute>
                            <Login />
                        </PublicRoute>
                    }
                />
                <Route
                    path="/register"
                    element={
                        <PublicRoute>
                            <Register />
                        </PublicRoute>
                    }
                />
                <Route
                    path="/reset-password/:hash"
                    element={
                        <PublicRoute>
                            <ChangePassword />
                        </PublicRoute>
                    }
                />
                <Route
                    path="/reset-password"
                    element={
                        <PublicRoute>
                            <ResetPassword />
                        </PublicRoute>
                    }
                />
                <Route
                    path="/confirm/:hash"
                    element={
                        <PublicRoute>
                            <Confirm />
                        </PublicRoute>
                    }
                />

                {/* inne trasy np. */}
                {/* <Route path="/contact" element={<PublicRoute><Contact /></PublicRoute>} /> */}

                {/* 404 dla niezalogowanego */}
                <Route path="*" element={<NotFound />} />
            </Route>

            {/* Protected (z layoutem) */}
            <Route
                path="/"
                element={
                    <ProtectedRoute>
                        <AuthenticatedLayout />
                    </ProtectedRoute>
                }
            >
                <Route path="dashboard" element={<Dashboard />} />
                <Route path="settings" element={<Settings />} />

                {/* Dodawaj więcej stron tutaj */}
                <Route path="/" element={<Navigate to="/dashboard" replace />} />

                {/* 404 dla zalogowanego */}
                <Route path="*" element={<NotFound />} />
            </Route>
            <Route path="/" element={<Navigate to="/dashboard" replace />} />

            {/* Domyślne przekierowanie */}
            <Route path="*" element={<Navigate to="/dashboard" replace />} />
        </Routes>
    );
}


export default App
