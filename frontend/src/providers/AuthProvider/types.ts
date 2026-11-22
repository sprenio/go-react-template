import type { UserType } from '@/domain';

export type MeResponse = {
  user: UserType;
};
export type AuthContextType = {
  setLoginUser: (user: UserType) => void;
  logout: (code: string) => void;
  fetchMe: () => Promise<void>;
  appUser: UserType | null;
  meInProgress: boolean;
};
