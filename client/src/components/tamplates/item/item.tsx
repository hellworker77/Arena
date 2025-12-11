import {BaseItem} from "types/game/game-item/baseItem.ts";
import {getItemImageById} from "types/game/game-item/getItemImageById.ts";
import {Tooltip} from "@/components/tamplates/tooltip/tooltip.tsx";
import {ExocetTitle} from "@/components/tamplates/exocetTitle/exocetTitle.tsx";
import {getTextColor} from "types/game/game-item/getTextColor.ts";

interface ItemProps {
    item: BaseItem;
}

export const Item = ({item}: ItemProps) => {

    return (
        <Tooltip>
            <Tooltip.Trigger>
                <img
                    src={getItemImageById(item.itemID)}
                    alt={item.name}
                    className="w-full h-full object-contain p-1"
                />
            </Tooltip.Trigger>
            <Tooltip.Content className="bg-[rgba(0,0,0,0.7)] p-1 text-white">
                <Tooltip.Title>
                    <span className="font-exocet">{item.name}</span>

                </Tooltip.Title>
                <ExocetTitle color={getTextColor("item-property")}>Durability 20/21</ExocetTitle>
                <div className="font-exocet">Durability</div>
                <div>Defense: {item.defense}</div>
                <div>Level: {item.requiredLevel}</div>
                <Tooltip.Footer>Some footer if needed</Tooltip.Footer>
            </Tooltip.Content>
        </Tooltip>
    )
}