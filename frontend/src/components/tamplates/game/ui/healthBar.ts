export class HealthBar {
    private bar: Phaser.GameObjects.Graphics;
    private text: Phaser.GameObjects.Text;
    private readonly width: number;
    private readonly height: number;

    constructor(scene: Phaser.Scene, x: number, y: number, width: number, height: number, health: number, max: number) {
        this.width = width;
        this.height = height;
        this.bar = scene.add.graphics();
        this.text = scene.add.text(x, y - 10, "", { fontSize: "12px", color: "#fff" }).setOrigin(0.5);
        this.updatePosition(x, y);
        this.updateHealth(health, max);
    }

    updateHealth(health: number, max: number) {
        const ratio = Phaser.Math.Clamp(health / max, 0, 1);
        const color = Phaser.Display.Color.Interpolate.ColorWithColor(
            new Phaser.Display.Color(255, 0, 0),
            new Phaser.Display.Color(0, 255, 0),
            1,
            ratio
        );

        this.bar.clear();
        this.bar.fillStyle(Phaser.Display.Color.GetColor(color.r, color.g, color.b), 1);
        this.bar.fillRect(-this.width / 2, 0, this.width * ratio, this.height);
        this.bar.lineStyle(1, 0xffffff).strokeRect(-this.width / 2, 0, this.width, this.height);
        this.text.setText(`${health}/${max}`);
    }

    updatePosition(x: number, y: number) {
        this.bar.setPosition(x, y);
        this.text.setPosition(x, y - 10);
    }

    destroy() {
        this.bar.destroy();
        this.text.destroy();
    }
}