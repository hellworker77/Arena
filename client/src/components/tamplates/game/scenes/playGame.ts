import Phaser from "phaser";
import {arenaWebRTCService} from "../../../../ws/arenaWebRTCService.ts";
import {ArenaStateSystem} from "../systems/arenaStateSystem.ts";
import {PlayerSystem} from "../systems/playerSystem.ts";
import {EnemySystem} from "../systems/enemySystem.ts";
import {ProjectileSystem} from "../systems/projectileSystem.ts";

export class PlayGame extends Phaser.Scene {
    private arena!: ArenaStateSystem;
    private player!: PlayerSystem;
    private enemies!: EnemySystem;
    private projectiles!: ProjectileSystem;

    constructor() {
        super({ key: "PlayGame" });
    }

    create() {


        this.input.mouse?.disableContextMenu();

        this.arena = new ArenaStateSystem();
        this.player = new PlayerSystem(this);
        this.enemies = new EnemySystem(this);
        this.projectiles = new ProjectileSystem(this);

        arenaWebRTCService.connect().then(() => {
            arenaWebRTCService.onUpdate((state) => {
                this.arena.updateFromServer(state);
            });
        });

        this.arena.onStateUpdated = (state) => {
            if (!this.player.getSprite()) {
                this.player.init(state.player.x, state.player.y);
            }
            this.enemies.sync(state.enemies);
            this.projectiles.sync(state.projectiles);
        };

        this.input.on("pointerdown", (pointer: { leftButtonDown: () => never; worldX: number; worldY: number; }) => {
            if (pointer.leftButtonDown()) {
                arenaWebRTCService.moveTo(pointer.worldX, pointer.worldY);
                this.player.getSprite()?.setData("targetX", pointer.worldX);
                this.player.getSprite()?.setData("targetY", pointer.worldY);
            }
        });
    }

    update(_time: number, delta: number) {
        const dt = delta / 1000;
        const state = this.arena.getInterpolatedState();
        if (!state) return;

        this.player.update(state.player, dt);
        this.enemies.update(state.enemies);
        this.projectiles.update(state.projectiles);
    }
}
