import type {AppRoute} from "./appRoute.ts";
import {Home} from "../components/pages/home/home.tsx";
import type {JSX} from "react";
import {Route} from "react-router-dom";
import {Arena} from "../components/pages/arena/arena.tsx";
import {Auction} from "../components/pages/auction/auction.tsx";
import {Friends} from "../components/pages/friends/friends.tsx";
import {LeaderBoard} from "../components/pages/leaderBoard/leaderBoard.tsx";
import {Market} from "../components/pages/market/market.tsx";
import {Settings} from "../components/pages/settings/settings.tsx";
import {Social} from "../components/pages/social/social.tsx";
import {Mail} from "../components/pages/mail/mail.tsx";
import {PatchNotes} from "../components/pages/patchNotes/patchNotes.tsx";
import {News} from "../components/pages/news/news.tsx";
import {SkillTree} from "../components/pages/character/skillTree/skillTree.tsx";
import {Inventory} from "../components/pages/character/inventory/inventory.tsx";

export const ARENA_ROUTE = "arena";
export const AUCTION_ROUTE = "auction";
export const CHARACTER_ROUTE = "character";
export const INVENTORY_SUBROUTE = "inventory";
export const SKILL_TREE_SUBROUTE = "skill-tree";
export const FRIENDS_ROUTE = "friends";
export const HOME_ROUTE = "home";
export const LEADERBOARD_ROUTE = "leaderboard";
export const MAIL_ROUTE = "mail";
export const MARKET_ROUTE = "market";
export const NEWS_ROUTE = "news";
export const PATCH_NOTES_ROUTE = "patch-notes";
export const SETTINGS_ROUTE = "settings";
export const SOCIAL_ROUTE = "social";

export const app_routes: AppRoute[] = [
    {
        path: ARENA_ROUTE,
        component: <Arena/>,
        name: "arena",
        authorizationRequired: true,
        container: "both",
    },
    {
        path: AUCTION_ROUTE,
        component: <Auction/>,
        name: "auction",
        authorizationRequired: true,
        container: "side"
    },
    {
        path: CHARACTER_ROUTE,
        component: null,
        name: "character",
        authorizationRequired: true,
        container: "side",
        subroutes: [
            {
                path: INVENTORY_SUBROUTE,
                component: <Inventory/>,
                name: "inventory",
                authorizationRequired: true,
                container: "side"
            },
            {
                path: SKILL_TREE_SUBROUTE,
                component: <SkillTree/>,
                name: "skill-tree",
                authorizationRequired: true,
                container: "side"
            }
        ],
    },
    {
        path: FRIENDS_ROUTE,
        component: <Friends/>,
        name: "friends",
        authorizationRequired: true,
        container: "side"
    },
    {
        path: HOME_ROUTE,
        component: <Home/>,
        name: "home",
        container: "both"
    },
    {
        path: LEADERBOARD_ROUTE,
        component: <LeaderBoard/>,
        name: "leaderboard",
        container: "side"
    },
    {
        path: MARKET_ROUTE,
        component: <Market/>,
        name: "market",
        authorizationRequired: true,
        container: "side"
    },
    {
        path: NEWS_ROUTE,
        component: <News/>,
        name: "news",
        container: "side"
    },
    {
        path: PATCH_NOTES_ROUTE,
        component: <PatchNotes/>,
        name: "patch-notes",
        container: "side"
    },
    {
        path: SETTINGS_ROUTE,
        component: <Settings/>,
        name: "settings",
        container: "side"
    },
    {
        path: SOCIAL_ROUTE,
        component: <Social/>,
        name: "social",
        container: "side"
    },
    {
        path: MAIL_ROUTE,
        component: <Mail/>,
        name: "mail",
        authorizationRequired: true,
        container: "both"
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