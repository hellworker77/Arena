import {BaseItem} from "../../../@types/game/game-item/baseItem.ts";
import {Item} from "../item/item.tsx";

interface InventorySlotProps {
    item?: BaseItem;
}

export const InventorySlot = ({item}: InventorySlotProps) => {

    return (
        <div className="w-full h-full bg-[#282828] border border-[#3E3E3E] flex items-center justify-center">
            {item && <Item item={item} />}
        </div>
    )
}