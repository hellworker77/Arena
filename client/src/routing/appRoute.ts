import type {JSX} from "react";

export type AppRouteContainer = "side" | "nav" | "both";

export type AppRoute = {

    path: string;

    name: string;

    component: JSX.Element | null;

    subroutes?: AppRoute[];

    authorizationRequired?: boolean;

    container: AppRouteContainer;
}

