import {app_routes} from "../../../routing/routes.tsx";
import {useTranslation} from "react-i18next";
import {SideBarLink} from "./link/sideBarLink.tsx";
import {useCallback, useMemo} from "react";
import type {AppRoute} from "../../../routing/appRoute.ts";
import {useRedirectToFirstChild} from "../../../hooks/useRedirectToFirstChild.ts";

interface VisibleRoute extends AppRoute {
    level: number;
    children: VisibleRoute[];
}

export const Side = () => {
    const {t} = useTranslation();
    const isAuthorized = true; // TODO: replace with real authorization check
    useRedirectToFirstChild(isAuthorized);

    const segments = useMemo(() => location.pathname.split("/").filter(Boolean), [])

    const filteredRoutes = useMemo(() => {
        return app_routes
            .filter(route => route.container === "side" || route.container === "both")
            .filter(route => !route.authorizationRequired || isAuthorized );
    }, [isAuthorized])

    const getVisibleRoutes = useCallback(
        (routes: AppRoute[], level: number = 0, parentPath = ""): VisibleRoute[] => {
            return routes.map(route => {
                const fullPath = parentPath + "/" + route.path.replace(/^\/+/, ""); // создаем полный путь
                const routeSegments = route.path.split("/").filter(Boolean);
                const isActive = segments.includes(routeSegments[0]);

                const children =
                    (isActive || level === 0) && route.subroutes
                        ? getVisibleRoutes(route.subroutes, level + 1, fullPath)
                        : [];

                return {
                    ...route,
                    path: fullPath,
                    level,
                    children,
                };
            });
        },
        [segments]
    );

    const visibleRoutes: VisibleRoute[] = useMemo(
        () => getVisibleRoutes(filteredRoutes),
        [filteredRoutes, getVisibleRoutes]
    );

    const renderLinks = (routes: VisibleRoute[]) =>
        routes.map(route => (
            <div key={route.path}>
                <SideBarLink to={route.path} style={{ paddingLeft: `${route.level * 16}px` }}>
                    {t(route.name)}
                </SideBarLink>
                {route.children.length > 0 && renderLinks(route.children)}
            </div>
        ));

    return <aside className="w-64 h-full bg-indigo-950/5 p-4 select-none">{renderLinks(visibleRoutes)}</aside>;
};