import type {UserType} from '@/types';


export type MeResponse = {
    user: UserType
}
export type AuthContextType = {
    setLoginUser: (user: UserType) => void
    logout: (code:string) => void
    appUser: UserType|null
    meInProgress: boolean
}
