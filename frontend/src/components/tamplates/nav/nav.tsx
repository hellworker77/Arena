import {app_routes} from "../../../routing/routes.tsx";
import {NavBarLink} from "./link/navBarLink.tsx";
import {useTranslation} from "react-i18next";

export const Nav = () => {
    const {t} = useTranslation();

    const isAuthorized = true; // TODO: replace with real authorization check

    return (
        <header className="flex items-center justify-between px-6 py-4 h-16
    bg-black/20 shadow-lg border-b border-black/20 select-nonerelative">

            <div className="absolute left-1/2 transform -transient-x-1/2 flex space-x-6">
                {app_routes
                    .filter(route => route.container === "nav" || route.container === "both")
                    .filter(route => !route.authorizationRequired || isAuthorized)
                    .map((route, i) => (
                        <NavBarLink key={`route-${i}`} variant="light" to={route.path}>
                            {t(`${route.name}`)}
                        </NavBarLink>
                    ))}
            </div>
        </header>
    );
};