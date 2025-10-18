export class ParticleSystem {
    private readonly scene: Phaser.Scene;

    constructor(scene: Phaser.Scene) {
        this.scene = scene;
    }

    private createEmitter(key: string, x: number, y: number, config?: Phaser.Types.GameObjects.Particles.ParticleEmitterConfig): Phaser.GameObjects.Particles.ParticleEmitter {
        return this.scene.add.particles(x, y, key, config);
    }

    playImpactEffect(x: number, y: number): void {
        const emitter = this.createEmitter("particle", x, y, {
            speed: 100,
            lifespan: 2000,
            scale: { start: 0.5, end: 0 },
            alpha: { start: 1, end: 0 },
            quantity: 10
        });
        emitter.explode(10, x, y);
        this.scene.time.delayedCall(500, () => emitter.destroy());
    }

    playDeathEffect(x: number, y: number): void {
        const emitter = this.createEmitter("particle", x, y, {
            speed: 100,
            lifespan: 2000,
            scale: { start: 0.5, end: 0 },
            alpha: { start: 1, end: 0 },
            quantity: 10
        });
        emitter.explode(12, x, y);
        this.scene.time.delayedCall(500, () => emitter.destroy());
    }

    playProjectileTrail(follow: Phaser.GameObjects.Sprite): Phaser.GameObjects.Particles.ParticleEmitter {
        const emitter = this.createEmitter("particle", 0, 0, {
            speed: 0,
            lifespan: 300,
            scale: { start: 0.2, end: 0 },
            alpha: { start: 0.5, end: 0 },
            tint: 0xffffff,
            quantity: 1
        });
        emitter.startFollow(follow);
        return emitter;
    }
}