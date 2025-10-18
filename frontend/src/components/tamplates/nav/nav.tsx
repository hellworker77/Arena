import {app_routes} from "../../../routing/routes.tsx";
import {NavBarLink} from "./link/navBarLink.tsx";
import {useTranslation} from "react-i18next";

export const Nav = () => {
    const {t} = useTranslation();

    const isAuthorized = true; // TODO: replace with real authorization check

    return (
        <header className="flex items-center justify-between px-6 py-4 bg-gradient-to-b from-black via-gray-900 to-black shadow-lg border-b border-red-800 select-none">
            <p className="text-2xl font-bold text-red-500 drop-shadow-[0_0_8px_rgba(255,0,0,0.7)] tracking-wider cursor-pointer">
                Wow Source 2.0
            </p>
            <div className="flex space-x-6">
                {app_routes
                    .filter(route => route.container === "nav" || route.container === "both")
                    .filter(route => !route.authorizationRequired || isAuthorized)
                    .map(route => (
                        <NavBarLink to={route.path}>
                            {t(`${route.name}`)}
                        </NavBarLink>
                    ))}
            </div>
        </header>
    );
};