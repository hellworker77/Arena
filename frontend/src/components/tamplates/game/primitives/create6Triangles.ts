export function create6Triangles(
    scene: Phaser.Scene,
    key: string,
    size: number = 64,
    color: number = 0x00ffff,
    glowRadius: number = 12,
    glowIntensity: number = 0.0125,
    baseAngleDeg: number = 30,
    tipDistance: number = 0
) {
    const texSize = size + glowRadius * 2;
    const cx = texSize / 2;
    const cy = texSize / 2;
    const g = scene.add.graphics();


    for (let r = glowRadius; r > 0; r--) {
        const alpha = glowIntensity * (r / glowRadius);
        g.fillStyle(color, alpha);
        g.fillCircle(cx, cy, size / 2 + r);
    }

    g.fillStyle(color, 1);

    const count = 6;
    const radius = size / 2;
    const span = Phaser.Math.DegToRad(baseAngleDeg) * 0.97;

    for (let i = 0; i < count; i++) {
        const ang = (i * Math.PI * 2) / count;

        const tipX = cx + tipDistance * Math.cos(ang);
        const tipY = cy + tipDistance * Math.sin(ang);

        const x1 = cx + radius * Math.cos(ang - span);
        const y1 = cy + radius * Math.sin(ang - span);
        const x2 = cx + radius * Math.cos(ang + span);
        const y2 = cy + radius * Math.sin(ang + span);

        g.fillTriangle(tipX, tipY, x1, y1, x2, y2);
    }

    g.setBlendMode(Phaser.BlendModes.ADD);
    g.generateTexture(key, texSize, texSize);
    g.destroy();
}