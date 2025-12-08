import {Rarity} from "types/game/game-item/rarity.ts";

type textColorOptions = Record<string, string>;

const textColors: textColorOptions = {
    [Rarity.CRAP]: "#777777",
    [Rarity.COMMON]: "#FFFFFF",
    [Rarity.MAGIC]: "#6969FF",
    [Rarity.RARE]: "#FFD33E",
    [Rarity.EPIC]: "#8A2BE2",
    [Rarity.LEGENDARY]: "#C79045",
    [Rarity.MYTHICAL]: "#FF6969",
    [Rarity.SET]: "#00D318",
    [Rarity.UNIQUE]: "#C79045",
    [Rarity.CRAFTED]: "#FF8718",
    ["item-property"]: "#6969FF",
};

export const getTextColor = (v: Rarity | "item-property"): string => {
    return textColors[v] || "#FFFFFF";
}