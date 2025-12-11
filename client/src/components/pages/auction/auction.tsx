import {useApi} from "@/api/useApi.ts";
import {useEffect} from "react";

export const Auction = () => {

    const {data, error, loading, request } = useApi()

    useEffect(() => {
        request("/api/v1/test/key").then();
    }, []);

    return (
        <div>
            {data}
        </div>
    )
}