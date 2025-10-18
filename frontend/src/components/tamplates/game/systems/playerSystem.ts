import {PredictionSystem} from "./predictionSystem.ts";
import type {ArenaState} from "../../../../@types/game/arenaState.ts";

export class PlayerSystem {
    private readonly scene: Phaser.Scene;
    private sprite!: Phaser.GameObjects.Sprite;

    constructor(scene: Phaser.Scene) {
        this.scene = scene;
    }

    init(x: number, y: number) {
        this.sprite = this.scene.add.sprite(x, y, "player");
        this.sprite.setOrigin(0.5);

        this.sprite.postFX.addGlow(
            0x00ffff,       // цвет свечения, как в Geometry Arena
            2.0,      // outerStrength — интенсивность внешнего свечения
            0.5,      // innerStrength — внутреннее свечение (можно уменьшить или оставить 0)
            false,        // knockout — false, чтобы объект оставался видимым
            0.5,           // quality — качество свечения, 0.5–1 для баланса между качеством и FPS
        5                // distance — радиус свечения вокруг объекта);
        );
    }

    update(statePlayer: ArenaState["player"], dt: number) {
        if (!this.sprite) return;

        this.sprite.x = statePlayer.x;
        this.sprite.y = statePlayer.y;

        const healthRatio = statePlayer.health / statePlayer.maxHealth;
        this.sprite.setTint(PredictionSystem.interpolateColor(healthRatio));

        const rotationSpeed = Math.PI / 2;
        this.sprite.rotation += rotationSpeed * dt;
    }

    getSprite() {
        return this.sprite;
    }
}