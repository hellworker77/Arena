import type {ArenaState} from "../../../../@types/game/arenaState.ts";

export class ProjectileSystem {
    private scene: Phaser.Scene;
    private projectiles = new Map<string, Phaser.GameObjects.Sprite>();

    constructor(scene: Phaser.Scene) {
        this.scene = scene;
    }

    sync(projectiles: ArenaState["projectiles"]) {
        projectiles.forEach(p => {
            if (!this.projectiles.has(p.id)) {
                const sprite = this.scene.add.sprite(p.x, p.y, "projectile");
                this.projectiles.set(p.id, sprite);
            }
        });

        this.projectiles.forEach((sprite, id) => {
            if (!projectiles.find(p => p.id === id)) {
                sprite.destroy();
                this.projectiles.delete(id);
            }
        });
    }

    update(projectiles: ArenaState["projectiles"]) {
        projectiles.forEach(p => {
            const sprite = this.projectiles.get(p.id);
            if (!sprite) return;

            sprite.x = p.x;
            sprite.y = p.y;
        });
    }
}