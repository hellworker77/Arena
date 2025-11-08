interface CurrencyDisplayProps {
    value: number;
}

export const CurrencyDisplay = ({value}: CurrencyDisplayProps) => {
    return (
        <div className="bg-[#222222] flex flex-row items-center gap-2 h-full w-full px-2">
            <div className="bg-white rounded-full w-4 h-4"></div>
            <div className="text-white text-sm border-l-2 border-white pl-2">{value}</div>
        </div>
    )

}