interface InventorySlotProps {
    item: object | null;
}

export const InventorySlot = ({item}: InventorySlotProps) => {

    return (
        <div className="bg-[#282828] w-full h-full border border-[#3E3E3E]" />
    )
}