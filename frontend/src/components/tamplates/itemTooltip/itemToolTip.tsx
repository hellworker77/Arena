import type {Item} from "../../../@types/game/game-item/item.ts";
import {useTranslation} from "react-i18next";

interface ItemToolTipProps {
    item: Item;
}

type Property = {
    key: string;
    value: number;
}

export const ItemToolTip = ({item}: ItemToolTipProps) => {

    const {t} = useTranslation();

    const propertyEntries: Property[] = Object.entries(item)
        .filter(([_, value]) => typeof value === "number" && value !== 0)
        .map(([key, value]) => ({ key, value: value as number }));

    return (
        <div
            className="bg-natural-900 text-white border border-yellow-900 rounded-md p-3 sgadow-lg max-w-xs text-sm select-none">
            <div className="mb-2 font-bold text-yellow-400">{item.name}</div>
            <div className="mb-2 text-gray-300 text-xs uppercase">{item.rarity}</div>

            {propertyEntries.map((property, index) => (
                <div key={index} className="flex justify-between mb-0.5 last:mb-0">
                    <span>{t(property.key)}</span>
                    <span className="font-semibold">{property.value}</span>
                </div>
            ))}
        </div>
    )
}