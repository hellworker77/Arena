import Phaser from "phaser";

export class PreloadAssets extends Phaser.Scene {

    constructor() {
        super({
            key: "PreloadAssets"
        });
    }

    preload(): void {
        this.load.image("enemy", "assets/sprites/enemy.png");
        this.load.image("player", "assets/sprites/player.png");
        this.load.image("projectile", "assets/sprites/projectile.png");
    }

    create() : void {

        this.scene.start('PlayGame');
    }
}