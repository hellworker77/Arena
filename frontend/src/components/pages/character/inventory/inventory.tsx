interface Item {
    id: number;
    name: string;
    icon: string; // url или иконка
}

const equipmentSlots = ["Head", "Chest", "Legs", "Weapon", "Shield"];

const mockItems: Item[] = Array.from({ length: 60 }, (_, i) => ({
    id: i,
    name: `Item ${i + 1}`,
    icon: "https://via.placeholder.com/40", // замените на свои иконки
}));

export const Inventory = () => {
    return (
        <div className="flex w-full max-w-4xl mx-auto mt-10 p-4 bg-gray-900 border-4 border-gray-700 rounded-lg text-white">
            {/* Экипировка */}
            <div className="w-1/3 flex flex-col items-center border-r-2 border-gray-700 pr-4">
                <h2 className="mb-4 text-xl font-bold">Экипировка</h2>
                {equipmentSlots.map((slot) => (
                    <div
                        key={slot}
                        className="w-20 h-20 mb-4 bg-gray-800 border-2 border-gray-600 flex items-center justify-center rounded-md hover:border-yellow-400 cursor-pointer"
                    >
                        {slot}
                    </div>
                ))}
            </div>

            {/* Сумка */}
            <div className="w-2/3 pl-4">
                <h2 className="mb-4 text-xl font-bold">Сумка</h2>
                <div className="h-[400px] overflow-y-auto border-2 border-gray-700 rounded-md p-2 bg-gray-800">
                    <div className="grid grid-cols-4 gap-2">
                        {mockItems.map((item) => (
                            <div
                                key={item.id}
                                className="w-20 h-20 bg-gray-700 border-2 border-gray-600 flex items-center justify-center rounded-md hover:border-yellow-400 cursor-pointer"
                            >
                                <img src={item.icon} alt={item.name} className="w-12 h-12" />
                            </div>
                        ))}
                    </div>
                </div>
            </div>
        </div>
    );
}