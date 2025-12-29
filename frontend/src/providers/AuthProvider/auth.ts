const tokenStorageKey = 'token';
const rememberMeStorageKey = 'rememberMe';

export function getRememberMeFlag():boolean {
    return localStorage.getItem(rememberMeStorageKey) === 'true'
}

export function saveToken(token:string, rememberMe = false) {
    if (rememberMe) {
        setLocalStorageToken(token)
    } else {
        setSessionStorageToken(token)
    }
}

export function getToken():string | null {
    return localStorage.getItem(tokenStorageKey) || sessionStorage.getItem(tokenStorageKey)
}

export function clearToken():void {
    clearLocalStorageToken()
    clearSessionToken()
}

export function isLoggedIn():boolean {
    return !!getToken()
}

function setLocalStorageToken(token:string) {
    localStorage.setItem(tokenStorageKey, token)
    clearSessionToken()
    setRememberMeFlag(true)
}
function setSessionStorageToken(token:string) {
    sessionStorage.setItem(tokenStorageKey, token)
    clearLocalStorageToken()
    setRememberMeFlag(false)
}
function clearSessionToken() {
    sessionStorage.removeItem(tokenStorageKey)
}
function clearLocalStorageToken() {
    localStorage.removeItem(tokenStorageKey)
}
function setRememberMeFlag(flagValue:boolean) {
    if (flagValue) {
        localStorage.setItem(rememberMeStorageKey, 'true')
    } else {
        localStorage.removeItem(rememberMeStorageKey)
    }
}