import Phaser from "phaser";
import {createSquare} from "../primitives/createSquare.ts";

export class PreloadAssets extends Phaser.Scene {

    constructor() {
        super({
            key: "PreloadAssets"
        });
    }

    preload(): void {
        this.load.image("enemy", "assets/sprites/enemy.png");

        createSquare(this, "player")

        this.load.image("particle", "assets/sprites/projectile.png");
        this.load.image("projectile", "assets/sprites/projectile.png");
    }

    create() : void {

        this.scene.start('PlayGame');
    }
}