import { AxiosError } from "axios";

export const getAPIErrorMessage = (error: unknown): string => {
    if (error instanceof AxiosError) {
        const serverMessage = error?.response?.data?.message || error?.response?.data?.error;

        if (serverMessage) return serverMessage;

        switch (error.response?.status) {
            case 400: return "Invalid request";
            case 401: return "Unauthorized";
            case 403: return "Forbidden";
            case 500: return "Server Error";
        }

        if (error.code === 'ERR_NETWORK') return "Network error. Please check your internet connection.";
        if (error.code === 'ECONNABORTED') return "The request timed out. Please try again.";

        return error.message || "Unexpected error";
    }

    return error instanceof Error ? error.message : "Unknown error";
}