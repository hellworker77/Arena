import {HealthBar} from "../ui/healthBar.ts";
import type {ArenaState} from "../../../../@types/game/arenaState.ts";

export class EnemySystem {
    private readonly scene: Phaser.Scene;
    private enemies = new Map<string, Phaser.GameObjects.Sprite>();
    private healthBars = new Map<string, HealthBar>();

    constructor(scene: Phaser.Scene) {
        this.scene = scene;
    }

    sync(enemies: ArenaState["enemies"]) {
        enemies.forEach(e => {
            if (!this.enemies.has(e.id)) {
                const sprite = this.scene.add.sprite(e.x, e.y, "enemy");
                this.enemies.set(e.id, sprite);
                this.healthBars.set(e.id, new HealthBar(this.scene, e.x, e.y - 30, 40, 5, e.health, e.maxHealth));
            }
        });

        this.enemies.forEach((sprite, id) => {
            if (!enemies.find(e => e.id === id)) {
                sprite.destroy();
                this.healthBars.get(id)?.destroy();
                this.healthBars.delete(id);
                this.enemies.delete(id);
            }
        });
    }

    update(enemies: ArenaState["enemies"]) {
        enemies.forEach(e => {
            const sprite = this.enemies.get(e.id);
            if (!sprite) return;

            sprite.x = e.x;
            sprite.y = e.y;

            const healthRatio = e.health / e.maxHealth;
            sprite.setTint(Phaser.Display.Color.GetColor(255 * (1 - healthRatio), 0, 255 * healthRatio));
            this.healthBars.get(e.id)?.updatePosition(sprite.x, sprite.y - 30);
            this.healthBars.get(e.id)?.updateHealth(e.health, e.maxHealth);
        });
    }
}