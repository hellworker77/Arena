import "./index.css";
import "./i18n/i18n.js";
import {createRoot} from "react-dom/client";
import {StrictMode, useEffect, useState} from "react";

const SplashScreen = () => {
    const [progress, setProgress] = useState(0);

    useEffect(() => {
        const timer = setInterval(() => {
            setProgress((p) => {
                if (p >= 100) {
                    clearInterval(timer);
                    return 100;
                }
                return p + 2;
            });
        }, 50);

        return () => clearInterval(timer);
    }, []);

    return (
        <div className="flex flex-col items-center justify-center h-screen bg-gray-900 text-white">
            <div className="flex flex-col items-center gap-4">
                <img
                    src="assets/logo.png"
                    alt="App Logo"
                    className="h-24 animate-pulse"
                />
                <h1 className="text-2xl font-semibold">Загрузка приложения...</h1>
            </div>

            <div className="w-64 h-3 bg-gray-700 rounded-full overflow-hidden mt-6">
                <div
                    className="bg-blue-500 h-full transition-all duration-300"
                    style={{ width: `${progress}%` }}
                />
            </div>

            <p className="mt-2 text-sm text-gray-400">{progress}%</p>
        </div>
    );
};

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <SplashScreen />
    </StrictMode>,
)