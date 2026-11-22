import { Outlet } from 'react-router-dom';
import { Footer } from '@/layouts/Footer';
import { Header } from '@/layouts/Header';
import { StateMessage } from '@/components/StateMessage';
import { Wrapper } from '@/layouts/Wrapper';

const PublicLayout = () => (
  <Wrapper>
    <Header />
    <StateMessage />
    <main className={`flex flex-1 items-center justify-center bg-[var(--bg)] p-[var(--space-xl)]`}>
      <Outlet />
    </main>
    <Footer />
  </Wrapper>
);

export default PublicLayout;
