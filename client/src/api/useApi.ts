import {useCallback, useState} from "react";
import {api} from "@/api/apiClient.ts";

export function useApi<T = any>() {
    const [data, setData] = useState<T | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const request = useCallback(async (url: string, params?: any) => {
        setLoading(true);
        setError(null);

        try {
            const { data } = await api.get<T>(url, { params });
            setData(data);
            return data;
        } catch (e: any) {
            setError(e?.response?.data?.message || "Unknown error");
            throw e;
        } finally {
            setLoading(false);
        }
    }, [])

    return { data, error, loading, request };
}