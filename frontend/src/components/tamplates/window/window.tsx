import React, { useEffect, useRef, useState } from "react";

type Direction = "right" | "bottom" | "corner" | null;

interface WindowProps {
    children?: React.ReactNode;
    title?: string;
}

export const Window = ({ children, title = "Window" }: WindowProps) => {
    const [state, setState] = useState({
        x: 100,
        y: 100,
        width: 400,
        height: 300,
    });
    const [resizing, setResizing] = useState<Direction>(null);

    const windowRef = useRef<HTMLDivElement>(null);
    const draggingRef = useRef(false);
    const rel = useRef({ x: 0, y: 0 });
    const startPos = useRef({ x: 0, y: 0 });
    const startSize = useRef({ width: 0, height: 0 });
    const frame = useRef<number>(0);

    // ⏳ Авторазмер по контенту один раз — при монтировании
    useEffect(() => {
        const el = windowRef.current;
        if (!el) return;

        const content = el.querySelector(".window-content") as HTMLDivElement;
        if (!content) return;

        // подгоняем окно под контент + header
        const rect = content.getBoundingClientRect();
        setState(prev => ({
            ...prev,
            width: Math.max(rect.width + 40, 300),
            height: Math.max(rect.height + 60, 200),
        }));
    }, [children]);

    const startDrag = (e: React.MouseEvent<HTMLDivElement>) => {
        draggingRef.current = true;
        rel.current = { x: e.clientX - state.x, y: e.clientY - state.y };
    };

    const startResize = (dir: Direction) => (e: React.MouseEvent) => {
        setResizing(dir);
        startPos.current = { x: e.clientX, y: e.clientY };
        startSize.current = { width: state.width, height: state.height };
    };

    useEffect(() => {
        let { x, y, width, height } = state;
        const el = windowRef.current;
        if (!el) return;

        const onMouseMove = (e: MouseEvent) => {
            if (draggingRef.current) {
                x = e.clientX - rel.current.x;
                y = e.clientY - rel.current.y;
                cancelAnimationFrame(frame.current!);
                frame.current = requestAnimationFrame(() => {
                    el.style.left = `${x}px`;
                    el.style.top = `${y}px`;
                });
            }

            if (resizing) {
                const deltaX = e.clientX - startPos.current.x;
                const deltaY = e.clientY - startPos.current.y;

                const smoothedX = deltaX * (1 - Math.exp(-Math.abs(deltaX) / 95));
                const smoothedY = deltaY * (1 - Math.exp(-Math.abs(deltaY) / 53));

                if (resizing === "right" || resizing === "corner") {
                    width = Math.max(200, startSize.current.width + smoothedX);
                }
                if (resizing === "bottom" || resizing === "corner") {
                    height = Math.max(150, startSize.current.height + smoothedY);
                }

                cancelAnimationFrame(frame.current!);
                frame.current = requestAnimationFrame(() => {
                    el.style.width = `${width}px`;
                    el.style.height = `${height}px`;
                });
            }
        };

        const onMouseUp = () => {
            if (draggingRef.current || resizing) {
                setState({ x, y, width, height });
            }
            draggingRef.current = false;
            setResizing(null);
        };

        window.addEventListener("mousemove", onMouseMove);
        window.addEventListener("mouseup", onMouseUp);
        return () => {
            window.removeEventListener("mousemove", onMouseMove);
            window.removeEventListener("mouseup", onMouseUp);
        };
    }, [resizing]);

    return (
        <div
            ref={windowRef}
            className="absolute border border-gray-600 shadow-2xl rounded-md bg-white flex flex-col select-none"
            style={{
                left: state.x,
                top: state.y,
                width: state.width,
                height: state.height,
            }}
        >
            {/* Верхняя панель */}
            <div
                onMouseDown={startDrag}
                className="h-8 bg-blue-600 text-white flex items-center px-3 rounded-t-md cursor-grab active:cursor-grabbing"
            >
                <span className="text-sm font-semibold truncate">{title}</span>
            </div>

            {/* Контент */}
            <div className="window-content flex-1 overflow-auto text-sm text-gray-800 p-3">
                {children || (
                    <div className="space-y-2">
                        <p>🪟 Пример окна в стиле Windows</p>
                        <p>Содержимое подстраивает размер при открытии.</p>
                        <p>Дальше можно тянуть мышкой как настоящее окно.</p>
                    </div>
                )}
            </div>

            {/* Ресайзеры */}
            <div
                onMouseDown={startResize("right")}
                className="absolute top-0 right-0 h-full w-2 cursor-ew-resize"
            />
            <div
                onMouseDown={startResize("bottom")}
                className="absolute bottom-0 left-0 w-full h-2 cursor-ns-resize"
            />
            <div
                onMouseDown={startResize("corner")}
                className="absolute right-0 bottom-0 w-3 h-3 cursor-nwse-resize"
            />
        </div>
    );
};
