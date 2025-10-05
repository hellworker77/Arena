import type {JSX} from "react";

export type AppRoute = {

    path: string;

    name: string;

    component: JSX.Element | null;

    subroutes?: AppRoute[];
}

