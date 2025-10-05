export const Nav = () => {
    return (
        <header className="flex items-center justify-between px-6 py-4 bg-gradient-to-b from-black via-gray-900 to-black shadow-lg border-b border-red-800 select-none">
            <p className="text-2xl font-bold text-red-500 drop-shadow-[0_0_8px_rgba(255,0,0,0.7)] tracking-wider cursor-pointer">
                Wow Source 2.0
            </p>
            <div className="flex space-x-6">
                <a
                    href="#"
                    className="text-gray-300 hover:text-red-400 cursor-pointer transition drop-shadow-[0_0_6px_rgba(255,50,50,0.7)]"
                >
                    Home
                </a>
                <a
                    href="#"
                    className="text-gray-300 hover:text-red-400 cursor-pointer transition drop-shadow-[0_0_6px_rgba(255,50,50,0.7)]"
                >
                    About
                </a>
                <a
                    href="#"
                    className="text-gray-300 hover:text-red-400 cursor-pointer transition drop-shadow-[0_0_6px_rgba(255,50,50,0.7)]"
                >
                    Contact
                </a>
            </div>
        </header>
    );
};