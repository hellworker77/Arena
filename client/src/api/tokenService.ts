export const tokenService = {
    getAccessToken: () => localStorage.getItem("access_token"),
    getRefreshToken: () => localStorage.getItem("refresh_token"),

    setAccessToken: (accessToken: string) => localStorage.setItem("access_token", accessToken),
    setRefreshToken: (refreshToken: string) => localStorage.setItem("refresh_token", refreshToken),



    setTokens: (accessToken: string, refreshToken?: string) => {
        tokenService.setAccessToken(accessToken);
        if(refreshToken) tokenService.setRefreshToken(refreshToken);
    },

    clearTokens: () => {
        localStorage.removeItem("access_token");
        localStorage.removeItem("refresh_token");
    }
}