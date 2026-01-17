import { Route, Routes} from 'react-router-dom';
import { getPublicRoutes } from '@/routes/PublicRoute';
import { getProtectedRoutes } from '@/routes/ProtectedRoute';

import { useAuth } from '@/providers/AuthProvider';
import Loader from '@/providers/LoaderProvider/Loader';
import {ProtectedRoute} from '@/routes/ProtectedRoute';
import {PublicRoute} from '@/routes/PublicRoute';
import {NotFound} from '@/pages/NotFound';

function App() {
    const { appUser, meInProgress } = useAuth();

    if (meInProgress) {
        return <Loader/>;
    }

    return (
        <Routes>
            {appUser ? getProtectedRoutes() : getPublicRoutes()}
            <Route element={appUser ? <ProtectedRoute /> : <PublicRoute/>}>
                <Route path="*" element={<NotFound />} />
            </Route>
        </Routes>
    );
}

export default App;