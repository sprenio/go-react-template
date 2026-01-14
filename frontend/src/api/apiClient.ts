import { logoutReasonCodes, apiErrorCodes } from './apiCodes';
import type {LogoutHandler, ApiFetchOptions, LoaderHandlers, ApiSuccessResponse, SuccessHandler} from './types';
import {ApiError} from './types';

let logoutHandler: LogoutHandler | null = null;
let loaderHandlers: LoaderHandlers | null = null;
let successHandler: SuccessHandler | null = null;

export function setLogoutHandler(fn:LogoutHandler | null) {
    logoutHandler = fn;
}

export function setLoaderHandlers(handlers: LoaderHandlers) {
    loaderHandlers = handlers;
}
export function setSuccessHandler(handler: SuccessHandler){
    successHandler = handler;
}
export async function apiFetch<T=unknown>(path:string, options:ApiFetchOptions = {}, showLoader = true, loaderText = ''):Promise<ApiSuccessResponse<T>> {

    const headers = {
        'Content-Type': 'application/json',
        'X-Frontend-Base-URL': window.location.origin,
        'credentials' : 'include',
        ...options.headers,
    };
    showLoader && loaderHandlers?.showLoader(loaderText);
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
                if (typeof logoutHandler === 'function') {
                    logoutHandler(logoutReasonCodes.SESSION_EXPIRED);
                }
                err.message = 'Unauthorized';
            }

            console.error(`API error ${res.status}; path: ${path}`, err);
            throw err;
        }

        if (isJson) {
            if (typeof successHandler === 'function') {
                successHandler();
            }
            return data;
        }

        console.warn('unexpected response type', contentType, data);
        throw new ApiError('unexpected response type', apiErrorCodes.UNKNOWN_ERROR);
    } catch (err:unknown) {
        if (err instanceof ApiError) {
            throw err;
        }
        if (err instanceof Error) {
            throw new ApiError(err.message, apiErrorCodes.UNKNOWN_ERROR);
        }
        throw new ApiError('Unknown error', apiErrorCodes.UNKNOWN_ERROR);
    } finally {
        showLoader && loaderHandlers?.hideLoader();
    }
}

export const api = {
    get: <T>(path:string, showLoader = true, loaderText = ''):Promise<ApiSuccessResponse<T>> => apiFetch(path, {}, showLoader, loaderText),
    post: <T>(path:string, body: any, loaderText = ''):Promise<ApiSuccessResponse<T>> =>
        apiFetch(path, { method: 'POST', body: JSON.stringify(body) }, true, loaderText),
    put: <T>(path:string, body: string, loaderText = ''):Promise<ApiSuccessResponse<T>> =>
        apiFetch(path, { method: 'PUT', body: JSON.stringify(body) }, true, loaderText),
    delete: <T>(path:string, loaderText = ''):Promise<ApiSuccessResponse<T>> => apiFetch(path, { method: 'DELETE' }, true, loaderText),
};
