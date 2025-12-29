import type {UserType} from '@/types';


export type MeResponse = {
    user: UserType
}
export type AuthContextType = {
    token: string|null
    setLoginUser: (user: UserType) => void
    setLoginToken: (token:string) => void
    logout: (code:string) => void
    user: UserType|null
    meInProgress: boolean
}
