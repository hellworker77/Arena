import type {AppRoute} from "./appRoute.ts";
import {Home} from "../components/pages/home/home.tsx";
import type {JSX} from "react";
import {Route} from "react-router-dom";
import {Scene} from "../components/pages/scene/scene.tsx";

export const HOME_ROUTE = "home";
export const SCENE_ROUTE = "scene";

export const app_routes: AppRoute[] = [
    {
        path: HOME_ROUTE,
        component: <Home/>,
        name: "home",
    },
    {
        path: SCENE_ROUTE,
        component: <Scene/>,
        name: "scene",
    }
]

export const renderRoutes = (
    routes: typeof app_routes,
    basePath: string = ""): JSX.Element[] => {
    return routes.flatMap(route => {
        const fullPath = `${basePath}/${route.path}`.replace(/\/+/g, "/");

        const currentRoute = <Route key={fullPath} path={fullPath} element={route.component}/>;

        if (route.subroutes) {
            return [currentRoute, ...renderRoutes(route.subroutes, fullPath)];
        }

        return [currentRoute];
    })
}