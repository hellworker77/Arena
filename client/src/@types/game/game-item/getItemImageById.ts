import shortsword_graphic from "../../../assets/itemGraphics/swords/shortsword_graphic.png";
import {ItemBase} from "@/assets/itemsStore/itemBase.ts";

const itemIdImageMap: Record<ItemBase, string> = {
    [ItemBase.SHORT_SWORD]: shortsword_graphic,
};

export const getItemImageById = (itemId: ItemBase) => {
    return itemIdImageMap[itemId] || undefined;
}