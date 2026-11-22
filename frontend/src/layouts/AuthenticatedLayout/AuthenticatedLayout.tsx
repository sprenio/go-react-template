import { Footer } from '@/layouts/Footer';
import { Header } from '@/layouts/Header';
import { StateMessage } from '@/components/StateMessage';
import { Wrapper } from '@/layouts/Wrapper';

export function AuthenticatedLayout({ children }: { children: React.ReactNode }) {
  return (
    <Wrapper>
      <Header />
      <StateMessage />
      <main
        className={'mx-auto w-full max-w-[1200px] items-start bg-[var(--bg)] p-[var(--space-xl)]'}
      >
        {children}
      </main>
      <Footer />
    </Wrapper>
  );
}

export default AuthenticatedLayout;
