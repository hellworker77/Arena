import axios from "axios";
import {tokenService} from "@/api/tokenService.ts";

export const api = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL,
    timeout: 5000,
})

let isRefreshing = false;
let queue: Array<(token: string) => void> = [];

api.interceptors.request.use((config) => {
    const token = tokenService.getAccessToken()

    if (token) {
        config.headers = config.headers || {};
        config.headers['Authorization'] = `Bearer ${token}`;
    }

    return config;
})

api.interceptors.response.use(
    (res) => res,
    async (error) => {
        const originalRequest = error.config;

        if (error?.response?.status === 401 && !originalRequest._retry) {
            if (isRefreshing) {
                return new Promise((resolve) => {
                    queue.push((newToken) => {
                        originalRequest.headers.Authorization = `Bearer ${newToken}`;
                        resolve(api(originalRequest));
                    })
                })
            }


            originalRequest._retry = true;
            isRefreshing = true;

            try {
                const refreshToken = tokenService.getRefreshToken();
                const {data} = await axios.post(`${import.meta.env.VITE_API_BASE_URL}/auth/refresh`, {
                    refreshToken
                });

                tokenService.setTokens(data.accessToken, data.refreshToken);

                queue.forEach((cb) => cb(data.accessToken));
                queue = [];

                isRefreshing = false;

                originalRequest.headers.Authorization = `Bearer ${data.accessToken}`;
                return api(originalRequest);
            } catch (err) {
                tokenService.clearTokens()
                isRefreshing = false;
                queue = []
                return Promise.reject(err);
            }
        }

        return Promise.reject(error);
    }
);