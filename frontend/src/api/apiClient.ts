import { clearToken, getToken } from '@/providers/AuthProvider';
import { logoutReasonCodes, apiErrorCodes } from './apiCodes';
import type {LogoutHandler, TokenRefresher, ApiFetchOptions, LoaderHandlers, ApiSuccessResponse} from './types';
import {ApiError} from './types';

let logoutHandler: LogoutHandler | null = null;
let tokenRefresher: TokenRefresher | null = null;
let loaderHandlers: LoaderHandlers | null = null;

export function setLogoutHandler(fn:LogoutHandler | null) {
    logoutHandler = fn;
}

export function setTokenRefresher(fn:TokenRefresher | null) {
    tokenRefresher = fn;
}
export function setLoaderHandlers(handlers: LoaderHandlers) {
    loaderHandlers = handlers;
}
export async function apiFetch<T=unknown>(path:string, options:ApiFetchOptions = {}, loaderText = ''):Promise<ApiSuccessResponse<T>> {
    const token = getToken();

    const headers = {
        'Content-Type': 'application/json',
        'X-Frontend-Base-URL': window.location.origin,
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
        ...options.headers,
    };
    loaderHandlers?.showLoader(loaderText);
    try {
        const res = await fetch(`/api/${path.replace(/^\/+/, '')}`, {
            ...options,
            headers,
        });

        const contentType = res.headers.get('content-type');
        const isJson = contentType?.includes('application/json');
        const data = isJson ? await res.json() : await res.text();

        if (!res.ok) {
            let apiErrorCode = apiErrorCodes.UNKNOWN_ERROR;
            if (Object.values(apiErrorCodes).includes(res.status)) {
                apiErrorCode = res.status;
            } else {
                console.warn(`Unexpected API error: ${res.status} ${res.statusText}`);
            }

            const err = new ApiError('', apiErrorCode);
            const defError = 'Server_error';
            if (isJson) {
                err.data = data;
                err.message = data?.message || defError;
                if (data?.code) {
                    err.code = data.code;
                }
            } else {
                err.message = data || defError;
            }

            if (err.code === apiErrorCodes.UNAUTHORIZED) {
                clearToken();
                if (typeof logoutHandler === 'function') {
                    logoutHandler(logoutReasonCodes.SESSION_EXPIRED);
                }
                err.message = 'Unauthorized';
            }

            console.error(`API error ${res.status}; path: ${path}`, err);
            throw err;
        }

        if (isJson) {
            // Refresh token obsłużony tutaj
            if (data?.token && typeof tokenRefresher === 'function') {
                tokenRefresher(data.token);
            }
            return data;
        }

        console.warn('unexpected response type', contentType, data);
        throw new ApiError('unexpected response type', apiErrorCodes.UNKNOWN_ERROR);
    } catch (err:unknown) {
        if (err instanceof ApiError) {
            console.log('instanceof apierror', err)
            throw err;
        }
        if (err instanceof Error) {
            throw new ApiError(err.message, apiErrorCodes.UNKNOWN_ERROR);
        }
        throw new ApiError('Unknown error', apiErrorCodes.UNKNOWN_ERROR);
    } finally {
        loaderHandlers?.hideLoader();
    }
}

export const api = {
    get: <T>(path:string, loaderText = ''):Promise<ApiSuccessResponse<T>> => apiFetch(path, {}, loaderText),
    post: <T>(path:string, body: any, loaderText = ''):Promise<ApiSuccessResponse<T>> =>
        apiFetch(path, { method: 'POST', body: JSON.stringify(body) }, loaderText),
    put: <T>(path:string, body: string, loaderText = ''):Promise<ApiSuccessResponse<T>> =>
        apiFetch(path, { method: 'PUT', body: JSON.stringify(body) }, loaderText),
    delete: <T>(path:string, loaderText = ''):Promise<ApiSuccessResponse<T>> => apiFetch(path, { method: 'DELETE' }, loaderText),
};
