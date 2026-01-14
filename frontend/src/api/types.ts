export type LogoutHandler = (reason: string) => void;
export class ApiError extends Error {
    code: number;
    data: unknown;

    constructor(message: string, code = 0, data: unknown = null) {
        // 1. Wywołujemy konstruktor Error
        super(message);

        // 2. Przypisujemy wartości
        this.code = code;
        this.data = data;

        // 3. Ustawiamy nazwę (ładniejsze logi w konsoli)
        this.name = 'ApiError';

        // 4. NAPRAWA PROTOTYPU (Kluczowe dla "instanceof")
        // Ustawiamy prototyp jawnie na ApiError
        Object.setPrototypeOf(this, ApiError.prototype);

        // 5. Opcjonalnie: zachowanie czytelnego stosu wywołań (V8/Chrome)
        if ((Error as any).captureStackTrace) {
            (Error as any).captureStackTrace(this, ApiError);
        }
    }
}
export type ApiFetchOptions = RequestInit & {
    headers?: Record<string, string>;
};
export type LoaderHandlers = {
    showLoader: (text?: string) => void
    hideLoader: () => void
}
export type SuccessHandler = () => void;

export interface ApiSuccessResponse<T = undefined> {
    code: number;
    message: string;
    data: T;
}

