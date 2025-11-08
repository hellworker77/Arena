import {useLocation, useNavigate} from "react-router-dom";
import {useEffect, useMemo} from "react";
import {app_routes} from "../routing/routes.tsx";
import type {AppRoute} from "../routing/appRoute.ts";

const findRouteBySegments = (routes: AppRoute[], segments: string[]): AppRoute | null => {
    if (segments.length === 0) return null;
    const [current, ...rest] = segments;
    const route = routes.find(r => r.path.replace(/^\/+/, "") === current) ?? null;
    if (!route) return null;
    if (rest.length === 0) return route;
    return route.subroutes ? findRouteBySegments(route.subroutes, rest) : route;
};

// рекурсивно ищем первый дочерний маршрут (или fallback на первый в списке)
const getRedirectPath = (route: AppRoute, parentPath = ""): string | null => {
    const pathPrefix = parentPath ? `${parentPath}/` : "";
    if (route.subroutes?.length) return `${pathPrefix}${route.subroutes[0].path}`;
    return `${pathPrefix}${route.path}`;
};

export const useRedirectToFirstChild = (isAuthorized: boolean) => {
    const location = useLocation();
    const navigate = useNavigate();

    const { routeToRedirect } = useMemo(() => {
        const segments = location.pathname.split("/").filter(Boolean);

        const authorizedRoutes = app_routes.filter(
            r => (r.container === "side" || r.container === "both") && (!r.authorizationRequired || isAuthorized)
        );

        const currentRoute = findRouteBySegments(authorizedRoutes, segments);

        if (currentRoute && !currentRoute.component) {
            const parentPath = segments.slice(0, -1).join("/");
            const redirectPath = getRedirectPath(currentRoute, parentPath);
            return { routeToRedirect: redirectPath };
        }

        return { routeToRedirect: null };
    }, [location.pathname, isAuthorized]);

    useEffect(() => {
        if (routeToRedirect) {
            navigate(`/${routeToRedirect}`, { replace: true });
        }
    }, [routeToRedirect, navigate]);
};