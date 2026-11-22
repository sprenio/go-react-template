import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
    return twMerge(clsx(inputs));
}

export function getDateMonthOffset(baseDate?: string, offset: number = 0) {
    const d = baseDate ? new Date(baseDate) : new Date();
    d.setMonth(d.getMonth() + offset);
    return d;
}

export function getDefaultMinDate(baseDate?: string) {
    return getDateMonthOffset(baseDate, -4);
}

export function formatShortDate(iso: string) {
    const date = new Date(iso);
    const day = String(date.getUTCDate()).padStart(2, '0');
    const month = String(date.getUTCMonth() + 1).padStart(2, '0');
    return `${day}.${month}`;
}