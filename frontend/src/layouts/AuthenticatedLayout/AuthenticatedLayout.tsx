import { Outlet } from 'react-router-dom';
import {Footer} from '@/layouts/Footer';
import {Header} from '@/layouts/Header';
import {StateMessage} from '@/components/StateMessage';
import {Wrapper} from '@/layouts/Wrapper'

function AuthenticatedLayout() {
  return (
    <Wrapper>
        <Header />
        <StateMessage />
        <main className={'items-start w-full max-w-[1200px] mx-auto p-[var(--space-xl)] bg-[var(--bg)]'}>
          <Outlet />
        </main>
        <Footer />
    </Wrapper>
  );
}

export default AuthenticatedLayout;
