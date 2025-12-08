import {Page} from "../../tamplates/page/page.tsx";
import {useEffect, useRef} from "react";
import Phaser from "phaser";
import {gameOptions} from "../../tamplates/game/options/gameOptions.ts";
import {Boot} from "../../tamplates/game/scenes/boot.ts";
import {Play} from "../../tamplates/game/scenes/play.ts";

export const Arena = () => {

    const divRef = useRef<HTMLDivElement | null>(null);
    const gameRef = useRef<Phaser.Game | null>(null);

    useEffect(() => {
        const container = divRef.current;
        if(!container) return;

        const scaleObject = {
            mode: Phaser.Scale.FIT,
            autoCenter: Phaser.Scale.NONE,
            width: container.clientWidth,
            height: container.clientHeight,
            parent: container,
        }

        const configObject: Phaser.Types.Core.GameConfig = {
            type: Phaser.CANVAS,
            backgroundColor: gameOptions.gameBackgroundColor,

            render: {
                pixelArt: true,
                antialias: false
            },

            scale: scaleObject,

            scene: [Boot, Play],
            physics: {
                default: "arcade",
                arcade: {debug: false}
            }
        };

        gameRef.current = new Phaser.Game(configObject);
        const game = gameRef.current;
        if(!game) return;

        return () => {
            gameRef.current?.destroy(true);
        };
    }, []);

    return (
        <Page>
            <Page.Body>
                <div ref={divRef}
                     className="w-full h-full">
                </div>
            </Page.Body>
        </Page>
    )
}