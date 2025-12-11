import Phaser from "phaser";
import {createSquare} from "../primitives/createSquare.ts";

export class Boot extends Phaser.Scene {

    constructor() {
        super({
            key: "boot"
        });
    }

    preload(): void {

    }

    create() : void {
 /*       createSquare(this, "enemy", 0xff0000, 10);
        createSquare(this, "particle", 0x00ff00, 40);
        createSquare(this, "projectile", 0xff0000, 10);*/
        createSquare(this, "player",  40, 0x00ff00);
        this.scene.start("play");
    }
}