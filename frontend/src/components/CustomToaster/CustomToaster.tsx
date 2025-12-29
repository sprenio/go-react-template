import { Toaster } from 'react-hot-toast';

export function CustomToaster() {
  return (
    <Toaster
      position="top-right"
      containerStyle={{ top: '4rem' }}
      gutter={8}
      toastOptions={{
        duration: 4000,
        style: {
          background: '#1f2937',
          color: '#fff',
          marginTop: '0.2rem',
        },
        success: {
          iconTheme: {
            primary: '#10b981', // green-500
            secondary: '#fff',
          },
        },
      }}
    />
  );
}
