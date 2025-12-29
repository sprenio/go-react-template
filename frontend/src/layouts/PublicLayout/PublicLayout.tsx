import {Outlet} from 'react-router-dom';
import {Footer} from '@/layouts/Footer';
import {Header} from '@/layouts/Header';
import {StateMessage} from '@/components/StateMessage';
import {Wrapper} from '@/layouts/Wrapper'

const PublicLayout = () => (
    <Wrapper>
        <Header/>
        <StateMessage/>
        <main className={`flex-1 flex items-center justify-center p-[var(--space-xl)] bg-[var(--bg)]`}>
            <Outlet/>
        </main>
        <Footer/>
    </Wrapper>
);

export default PublicLayout;
