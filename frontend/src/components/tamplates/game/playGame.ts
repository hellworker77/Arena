import { gameOptions } from "./gameOptions.ts";

export class PlayGame extends Phaser.Scene {
    private playerGroup!: Phaser.Physics.Arcade.Group;
    private enemyGroup!: Phaser.Physics.Arcade.Group;
    private projectileGroup!: Phaser.Physics.Arcade.Group;

    private selectedUnits: Phaser.Types.Physics.Arcade.SpriteWithDynamicBody[] = [];
    private selectionRect!: Phaser.GameObjects.Graphics | null;
    private selectionStart!: Phaser.Math.Vector2 | null;

    constructor() {
        super({ key: "PlayGame" });
    }

    create(): void {
        // 🚫 отключаем контекстное меню (ПКМ)
        this.input.mouse?.disableContextMenu();

        // 🎯 группы
        this.playerGroup = this.physics.add.group();
        this.enemyGroup = this.physics.add.group();
        this.projectileGroup = this.physics.add.group();

        // 🧍 создаём несколько игроков
        for (let i = 0; i < 3; i++) {
            const unit = this.physics.add.sprite(400 + i * 40, 300, "player");
            unit.setData("config", { ...gameOptions.player });
            unit.setData("moveTarget", null);
            unit.setData("lastAttackTime", 0);
            this.playerGroup.add(unit);
        }

        // 🧟 враги спавнятся
        this.time.addEvent({
            delay: gameOptions.enemySpawnRate,
            loop: true,
            callback: () => {
                const spawn = Phaser.Geom.Rectangle.Random(
                    new Phaser.Geom.Rectangle(0, 0, gameOptions.gameSize.width, gameOptions.gameSize.height)
                );
                const enemy = this.physics.add.sprite(spawn.x, spawn.y, "enemy");
                enemy.setData("config", { ...gameOptions.enemy });
                enemy.setData("lastAttackTime", 0);
                this.enemyGroup.add(enemy);
            },
        });

        // 🖱 управление
        this.setupSelectionControls();

        // 💥 столкновения пуль и врагов
        this.physics.add.collider(this.projectileGroup, this.enemyGroup, (projObj, enemyObj) => {
            const projectile = projObj as Phaser.Types.Physics.Arcade.SpriteWithDynamicBody;
            const enemy = enemyObj as Phaser.Types.Physics.Arcade.SpriteWithDynamicBody;
            const dmg = projectile.getData("damage") ?? 10;

            enemy.setData("config", {
                ...enemy.getData("config"),
                health: enemy.getData("config").health - dmg,
            });

            if (enemy.getData("config").health <= 0) {
                this.enemyGroup.killAndHide(enemy);
                enemy.body.checkCollision.none = true;
            }

            this.projectileGroup.killAndHide(projectile);
            projectile.body.checkCollision.none = true;
        });
    }

    update(time: number, delta: number): void {
        // 🚶 движение игроков
        this.playerGroup.getChildren().forEach((obj) => {
            const unit = obj as Phaser.Types.Physics.Arcade.SpriteWithDynamicBody;
            const cfg = unit.getData("config");

            // 🧠 пропускаем мёртвых
            if (!unit.visible || cfg.health <= 0) return;

            const target: Phaser.Math.Vector2 | null = unit.getData("moveTarget");

            if (target) {
                const dist = Phaser.Math.Distance.Between(unit.x, unit.y, target.x, target.y);
                if (dist > 5) {
                    this.physics.moveTo(unit, target.x, target.y, cfg.speed);
                } else {
                    unit.body.setVelocity(0);
                    unit.setData("moveTarget", null);
                }
            } else {
                // если не движется — можно стрелять
                this.tryAttack(unit, time);
            }
        });

        // 🧟 движение врагов к ближайшему игроку
        this.enemyGroup.getChildren().forEach((obj) => {
            const enemy = obj as Phaser.Types.Physics.Arcade.SpriteWithDynamicBody;
            const cfg = enemy.getData("config");
            if (!enemy.visible || cfg.health <= 0) return;

            const alivePlayers = this.playerGroup.getChildren().filter((p) => {
                const pcfg = (p as any).getData("config");
                return (p as any).visible && pcfg.health > 0;
            });

            const closestPlayer = this.physics.closest(enemy, alivePlayers);
            if (closestPlayer) {
                const dist = Phaser.Math.Distance.Between(enemy.x, enemy.y, closestPlayer.x, closestPlayer.y);
                if (dist > cfg.attackRange) {
                    this.physics.moveToObject(enemy, closestPlayer, cfg.speed);
                } else {
                    enemy.body.setVelocity(0);
                    this.tryEnemyAttack(enemy, closestPlayer, time);
                }
            } else {
                // если игроков нет — стоять
                enemy.body.setVelocity(0);
            }
        });

        const aliveCount = this.playerGroup
            .getChildren()
            .filter((p) => (p as any).getData("config").health > 0).length;

        if (aliveCount === 0) {
            this.scene.restart();
        }
    }

    // 🔫 атака игрока
    private tryAttack(unit: Phaser.Types.Physics.Arcade.SpriteWithDynamicBody, now: number) {
        const cfg = unit.getData("config");
        if (cfg.health <= 0) return; // не стреляет мёртвый

        const lastAttack = unit.getData("lastAttackTime") || 0;
        if (now - lastAttack < cfg.attackSpeed) return;

        const enemies = this.enemyGroup.getChildren().filter((e) => {
            const ecfg = (e as any).getData("config");
            return (e as any).visible && ecfg.health > 0;
        });

        const closestEnemy = this.physics.closest(unit, enemies);
        if (!closestEnemy) return;

        const distance = Phaser.Math.Distance.Between(unit.x, unit.y, closestEnemy.x, closestEnemy.y);
        if (distance > cfg.attackRange) return;

        unit.setData("lastAttackTime", now);

        // создаём снаряд
        const projectile = this.physics.add.sprite(unit.x, unit.y, "projectile");
        const dmg = Phaser.Math.Between(cfg.damage.min, cfg.damage.max);
        projectile.setData("damage", dmg);
        this.projectileGroup.add(projectile);
        this.physics.moveToObject(projectile, closestEnemy, gameOptions.projectile.velocity);
    }

    // 💢 атака врага
    private tryEnemyAttack(
        enemy: Phaser.Types.Physics.Arcade.SpriteWithDynamicBody,
        target: Phaser.Types.Physics.Arcade.SpriteWithDynamicBody,
        now: number
    ) {
        const cfg = enemy.getData("config");
        const last = enemy.getData("lastAttackTime") || 0;

        if (now - last < cfg.attackSpeed) return;

        enemy.setData("lastAttackTime", now);

        const dmg = Phaser.Math.Between(cfg.damage.min, cfg.damage.max);
        const playerCfg = target.getData("config");
        playerCfg.health -= dmg;

        if (playerCfg.health <= 0) {
            this.playerGroup.killAndHide(target);
            target.body.checkCollision.none = true;
        }

        target.setData("config", playerCfg);
    }

    // 🖱 выделение / перемещение
    private setupSelectionControls() {
        this.selectionRect = this.add.graphics({ lineStyle: { color: 0x00ff00, width: 1 } });
        this.selectionStart = null;

        this.input.on("pointerdown", (pointer: Phaser.Input.Pointer) => {
            if (pointer.button === 0) {
                // левая кнопка
                const clicked = this.getUnitUnderPointer(pointer);
                if (clicked) {
                    this.clearSelection();
                    this.selectUnit(clicked);
                } else {
                    this.selectionStart = new Phaser.Math.Vector2(pointer.x, pointer.y);
                }
            }
            if (pointer.button === 2 && this.selectedUnits.length > 0) {
                this.moveSelectedUnits(pointer.worldX, pointer.worldY);
            }
        });

        this.input.on("pointerup", (pointer: Phaser.Input.Pointer) => {
            if (pointer.button === 0 && this.selectionStart) {
                const end = new Phaser.Math.Vector2(pointer.x, pointer.y);
                this.selectUnitsInRect(this.selectionStart, end);
                this.selectionRect!.clear();
                this.selectionStart = null;
            }
        });

        this.input.on("pointermove", (pointer: Phaser.Input.Pointer) => {
            if (this.selectionStart) {
                this.drawSelectionRect(this.selectionStart, new Phaser.Math.Vector2(pointer.x, pointer.y));
            }
        });
    }

    private getUnitUnderPointer(pointer: Phaser.Input.Pointer) {
        return this.playerGroup.getChildren().find((u) =>
            (u as Phaser.GameObjects.Sprite).getBounds().contains(pointer.x, pointer.y)
        ) as Phaser.Types.Physics.Arcade.SpriteWithDynamicBody | undefined;
    }

    private selectUnit(unit: Phaser.Types.Physics.Arcade.SpriteWithDynamicBody) {
        this.selectedUnits.push(unit);
        unit.setTint(0x00ff00);
    }

    private clearSelection() {
        this.selectedUnits.forEach((u) => u.clearTint());
        this.selectedUnits = [];
    }

    private selectUnitsInRect(start: Phaser.Math.Vector2, end: Phaser.Math.Vector2) {
        const rect = new Phaser.Geom.Rectangle(
            Math.min(start.x, end.x),
            Math.min(start.y, end.y),
            Math.abs(end.x - start.x),
            Math.abs(end.y - start.y)
        );
        this.clearSelection();

        this.playerGroup.getChildren().forEach((unitObj) => {
            const unit = unitObj as Phaser.Types.Physics.Arcade.SpriteWithDynamicBody;
            if (rect.contains(unit.x, unit.y)) {
                this.selectUnit(unit);
            }
        });
    }

    private drawSelectionRect(start: Phaser.Math.Vector2, end: Phaser.Math.Vector2) {
        this.selectionRect!.clear();
        this.selectionRect!.strokeRect(
            Math.min(start.x, end.x),
            Math.min(start.y, end.y),
            Math.abs(end.x - start.x),
            Math.abs(end.y - start.y)
        );
    }

    private moveSelectedUnits(x: number, y: number) {
        const formationSize = Math.ceil(Math.sqrt(this.selectedUnits.length));
        const spacing = 40;
        let index = 0;

        this.selectedUnits.forEach((unit) => {
            const row = Math.floor(index / formationSize);
            const col = index % formationSize;
            const tx = x + (col - formationSize / 2) * spacing;
            const ty = y + (row - formationSize / 2) * spacing;
            unit.setData("moveTarget", new Phaser.Math.Vector2(tx, ty));
            index++;
        });
    }
}
