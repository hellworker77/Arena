import {Inventory} from "./inventory.tsx";
import {Stash} from "./stash.tsx";

export const InventoryPage = () => {
    return (
        <div className="flex gap-4 h-full">
            <div className="flex-1 h-full">
                <Stash />
            </div>
            <div className="flex-1 h-full">
                <Inventory />
            </div>
        </div>
    )
}