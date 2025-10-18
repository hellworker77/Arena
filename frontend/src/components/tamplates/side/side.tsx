import {app_routes} from "../../../routing/routes.tsx";
import {useTranslation} from "react-i18next";
import {SideBarLink} from "./link/sideBarLink.tsx";

export const Side = () => {
    const {t} = useTranslation()

    const isAuthorized = true; // TODO: replace with real authorization check

    return (
        <aside className="w-64 h-full bg-gradient-to-b from-black via-gray-900 to-black border-r border-red-900 shadow-xl p-4 select-none">
            <nav className="space-y-3">
                {app_routes
                    .filter(route => !route.authorizationRequired || isAuthorized)
                    .filter(route => route.container === "side" || route.container === "both")
                    .map((subroutine) => (
                    <SideBarLink to={subroutine.path}>
                        {t(`${subroutine.name}`)}
                    </SideBarLink>
                ))}
            </nav>
        </aside>
    );
};