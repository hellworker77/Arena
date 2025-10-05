import { useEffect, useRef } from "react";
import * as PIXI from "pixi.js";

export const PixiCanvas = () => {
    const pixiContainer = useRef<HTMLDivElement | null>(null);
    const appRef = useRef<PIXI.Application | null>(null);
    const canvasRef = useRef<HTMLCanvasElement | null>(null);
    const mountedRef = useRef(false);

    useEffect(() => {
        if (mountedRef.current) return; // ⚡ гарантируем один запуск
        mountedRef.current = true;

        const container = pixiContainer.current;
        if (!container) return;

        const app = new PIXI.Application();
        appRef.current = app;

        (async () => {
            try {
                await app.init({
                    width: 800,
                    height: 600,
                    backgroundColor: 0x1099bb,
                    autoDensity: true,
                    resolution: window.devicePixelRatio || 1,
                });

                canvasRef.current = app.canvas;
                container.appendChild(app.canvas);

                const texture = await PIXI.Assets.load("https://pixijs.com/assets/bunny.png");
                const bunny = new PIXI.Sprite(texture);
                bunny.anchor.set(0.5);
                bunny.x = app.renderer.width / 2;
                bunny.y = app.renderer.height / 2;
                app.stage.addChild(bunny);

                app.ticker.add(() => {
                    bunny.rotation += 0.01;
                });
            } catch (err) {
                console.error("PixiJS init failed", err);
            }
        })();

        return () => {
            const app = appRef.current;
            if (app) {
                if (container && canvasRef.current && container.contains(canvasRef.current)) {
                    container.removeChild(canvasRef.current);
                }
                if (app.renderer) {
                    app.destroy(true, { children: true });
                }
                appRef.current = null;
                canvasRef.current = null;
            }
        };
    }, []);

    return <div ref={pixiContainer}></div>;
};
