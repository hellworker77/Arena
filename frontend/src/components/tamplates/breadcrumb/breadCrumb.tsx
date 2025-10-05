import React from "react";
import {app_routes} from "../../../routing/routes.tsx";
import {useEffect, useState} from "react";
import {NavLink, useLocation} from "react-router-dom";
import {useTranslation} from "react-i18next";

type Path = {
    name: string;

    path: string;
}

const findRouteChain = (
    routes: typeof app_routes,
    pathParts: string[],
    basePath = ""
): Path[] | null => {
    for (const route of routes) {
        const fullPath = `${basePath}/${route.path}`.replace(/\/+/g, "/");
        const match =
            route.path === pathParts[0] ||
            (route.path.startsWith(":") && pathParts.length > 0);

        if (match) {
            const currentPart = pathParts[0];
            const displayName = route.name.replace("{}", currentPart);
            const breadcrumb = { name: displayName, path: fullPath };

            if (pathParts.length === 1) {
                return [breadcrumb];
            }

            if (route.subroutes) {
                const subChain = findRouteChain(
                    route.subroutes,
                    pathParts.slice(1),
                    fullPath
                );

                if (subChain) {
                    return [breadcrumb, ...subChain];
                }
            }
        }
    }

    for (const route of routes) {
        if (route.subroutes) {
            const subChain = findRouteChain(route.subroutes, pathParts, `${basePath}/${route.path}`);
            if (subChain) {
                const displayName = route.name;
                const breadcrumb = {
                    name: displayName,
                    path: `${basePath}/${route.path}`.replace(/\/+/g, "/")
                };
                return [breadcrumb, ...subChain];
            }
        }
    }

    return null;
};

export const BreadCrumb = () => {
    const [breadCrumbs, setBreadCrumbs] = useState<Path[]>([]);

    const {t} = useTranslation();

    const location = useLocation();

    useEffect(() => {
        const pathParts = location.pathname.split("/").filter(Boolean);
        setBreadCrumbs(findRouteChain(app_routes, pathParts) ?? []);
    }, [location]);


    return (
        <div className="flex items-center space-x-2 text-red-500 font-serif tracking-wide select-none">
            {breadCrumbs.flatMap((crumb, index) => {
                const isLast = index === breadCrumbs.length - 1;

                return (
                    <React.Fragment key={crumb.path}>
                        {isLast ? (
                            <span className="text-red-400 drop-shadow-[0_0_6px_rgba(255,0,0,0.8)]">
              {t(crumb.name)}
            </span>
                        ) : (
                            <NavLink
                                to={crumb.path}
                                className="cursor-pointer hover:text-red-300 hover:drop-shadow-[0_0_6px_rgba(255,80,80,0.8)] transition"
                            >
                                {crumb.name}
                            </NavLink>
                        )}
                        {!isLast && (
                            <span className="text-gray-500">⟩</span>
                        )}
                    </React.Fragment>
                );
            })}
        </div>
    )
}