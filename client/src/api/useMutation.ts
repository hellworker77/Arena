import {useCallback, useState} from "react";
import {api} from "@/api/apiClient.ts";

type Method = "post" | "put" | "delete";

export function useMutation<T = any, B = any>(method: Method, url: string) {
    const [data, setData] = useState<T | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const mutate = useCallback(
        async (body?: B) => {
            setLoading(true);
            setError(null);
            try {
                const { data } = await api[method]<T>(url, body);
                setData(data);
                return data;
            } catch (err: any) {
                setError(err?.response?.data?.message || "Unknown error");
                throw err;
            } finally {
                setLoading(false);
            }
        },
        [method, url]
    );

    return { data, error, loading, mutate };
}