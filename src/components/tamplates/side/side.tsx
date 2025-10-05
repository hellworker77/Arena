export const Side = () => {
    return (
        <aside className="w-64 h-full bg-gradient-to-b from-black via-gray-900 to-black border-r border-red-900 shadow-xl p-4 select-none">
            <nav className="space-y-3">
                <a
                    href="#"
                    className="block px-3 py-2 rounded-lg text-gray-300 hover:text-red-400 hover:bg-gray-800 cursor-pointer transition drop-shadow-[0_0_6px_rgba(255,50,50,0.6)]"
                >
                    Dashboard
                </a>
                <a
                    href="#"
                    className="block px-3 py-2 rounded-lg text-gray-300 hover:text-red-400 hover:bg-gray-800 cursor-pointer transition drop-shadow-[0_0_6px_rgba(255,50,50,0.6)]"
                >
                    Characters
                </a>
                <a
                    href="#"
                    className="block px-3 py-2 rounded-lg text-gray-300 hover:text-red-400 hover:bg-gray-800 cursor-pointer transition drop-shadow-[0_0_6px_rgba(255,50,50,0.6)]"
                >
                    Inventory
                </a>
                <a
                    href="#"
                    className="block px-3 py-2 rounded-lg text-gray-300 hover:text-red-400 hover:bg-gray-800 cursor-pointer transition drop-shadow-[0_0_6px_rgba(255,50,50,0.6)]"
                >
                    Settings
                </a>
            </nav>
        </aside>
    );
};