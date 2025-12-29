import i18next from 'i18next';
import { toast } from 'react-hot-toast';
import {ApiError} from './types';

export const apiCodes = {
    API_General_Success: 1000,
    API_General_Unknown_Error: 1001,
    API_General_Invalid_JSON: 1002,

    API_General_Custom_Error: 1051,
    API_General_Invalid_Input_Value: 1052,

    API_Register_Success: 1100,
    API_Register_User_Name_Or_Email_Taken: 1101,

    API_Login_Invalid_Credentials: 1201,

    API_Confirm_Success: 1300,
    API_Settings_Success: 1400,
    API_Email_Change_Success: 1500,
    API_Reset_Password_Success: 1600,
    API_Password_Change_Success: 1700,
};

export const logoutReasonCodes = {
    SESSION_EXPIRED: 'SESSION_EXPIRED',
    USER_LOGGED_OUT: 'USER_LOGGED_OUT',
    TOKEN_INVALID: 'TOKEN_INVALID',
    UNKNOWN_ERROR: 'UNKNOWN_ERROR',
};
export const apiErrorCodes = {
    UNAUTHORIZED: 401,
    FORBIDDEN: 403,
    NOT_FOUND: 404,
    INTERNAL_SERVER_ERROR: 500,
    BAD_REQUEST: 400,
    CONFLICT: 409,
    UNPROCESSABLE_ENTITY: 422,
    SERVICE_UNAVAILABLE: 503,
    UNKNOWN_ERROR: 504,
};
export function toastErrorMessage(message:string) {
    toast.error(message);
}
export function toastError(error:ApiError) {
    toastErrorMessage(getApiCodeDescription(error.code));
}
export function toastSuccess(code:number) {
    let message = getApiCodeDescription(code);
    toast.success(message);
}

export function getApiCodeDescription(code:number) {
    const unknownCode = '__UNKNOWN__';
    const t = i18next.t;
    let description = t(`apiCodes.${code}`, { defaultValue: unknownCode });
    if (description === unknownCode) {
        console.warn(`No translation found for API code: ${code}`);
        let translation_code = code % 10 || code < 1000 ? 'unknown_error' : 'general_success';
        description = t('apiCodes.' + translation_code, { code: code });
    }
    return description;
}
