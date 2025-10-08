import {useEffect, useRef} from "react";
import Phaser from "phaser";
import {PreloadAssets} from "./preloadAssets";
import {PlayGame} from "./playGame";
import {gameOptions} from "./gameOptions";

export const Canvas = () => {
    const divRef = useRef<HTMLDivElement | null>(null);
    const gameRef = useRef<Phaser.Game | null>(null);

    useEffect(() => {
        const container = divRef.current;
        if(!container) return;

        const scaleObject = {
            mode: Phaser.Scale.NONE,
            autoCenter: Phaser.Scale.NONE,
            parent: container,
            width: gameOptions.gameSize.width,
            height: gameOptions.gameSize.height
        }

        const configObject: Phaser.Types.Core.GameConfig = {
            type: Phaser.WEBGL,
            backgroundColor: gameOptions.gameBackgroundColor,
            scale: scaleObject,
            scene: [
                PreloadAssets,
                PlayGame
            ],
            physics: {
                default: "arcade",
                arcade: {
                    debug: false
                }
            }
        };

        gameRef.current = new Phaser.Game(configObject);
        const game = gameRef.current;
        if(!game) return;

        return () => {
            gameRef.current?.destroy(true);
        };
    }, []);

    return <div ref={divRef}
                className="w-full h-full">
    </div>;
};