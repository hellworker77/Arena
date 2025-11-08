import {app_routes} from "../../../routing/routes.tsx";
import {useNavigate} from "react-router-dom";
import {useEffect} from "react";

type FoundRoute = { path: string; };

function findFirstNonNullRoute(routes: typeof app_routes, parentPath = ""): FoundRoute | null {
    for (const r of routes) {
        const fullPath = `${parentPath}/${r.path}`.replace(/\/+/g, "/");
        if (r.component) {
            return {path: fullPath};
        }
        if (r.subroutes) {
            const sub = findFirstNonNullRoute(r.subroutes, fullPath);
            if (sub) return sub;
        }
    }
    return null;
}

export const FirstRouteRedirectOrNotFound = () => {
    const navigate = useNavigate();
    const found = findFirstNonNullRoute(app_routes);

    useEffect(() => {
        if (found) {
            navigate(found.path, {replace: true});
        }
    }, [found, navigate]);

    return (
        <main className="flex items-center justify-center h-screen w-full bg-white text-gray-800">
            <div className="text-center px-4">
                <h1 className="text-6xl font-bold mb-4">404</h1>
                <p className="text-xl mb-6">Страница не найдена</p>
            </div>
        </main>
    );
};