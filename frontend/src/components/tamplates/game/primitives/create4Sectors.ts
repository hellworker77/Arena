export function create4Sectors(
    scene: Phaser.Scene,
    key: string,
    size: number = 64,
    color: number = 0xff00ff,
    glowRadius: number = 12,
    glowIntensity: number = 0.0125,
    gapDeg: number = 25,      // зазор между секторами (в градусах)
    innerRadius: number = 4  // расстояние от центра, с которого начинается сектор
) {
    const texSize = size + glowRadius * 2;
    const cx = texSize / 2;
    const cy = texSize / 2;
    const g = scene.add.graphics();

    // glow
    for (let r = glowRadius; r > 0; r--) {
        const alpha = glowIntensity * (r / glowRadius);
        g.fillStyle(color, alpha);
        g.fillCircle(cx, cy, size / 2 + r);
    }

    g.fillStyle(color, 1);

    const count = 4;
    const outerRadius = size / 2;
    const fullSector = 360 / count;
    const gapRad = Phaser.Math.DegToRad(gapDeg);
    const sectorSpan = Phaser.Math.DegToRad(fullSector) - gapRad;

    for (let i = 0; i < count; i++) {
        const startAngle = Phaser.Math.DegToRad(i * fullSector) + gapRad / 2;
        const endAngle = startAngle + sectorSpan;

        // Рисуем кольцевой сектор
        const points: Phaser.Types.Math.Vector2Like[] = [];

        // внешняя дуга
        const step = sectorSpan / 20; // сглаживание (чем больше — тем плавнее)
        for (let a = startAngle; a <= endAngle; a += step) {
            points.push({
                x: cx + outerRadius * Math.cos(a),
                y: cy + outerRadius * Math.sin(a),
            });
        }

        // внутренняя дуга (обратно)
        for (let a = endAngle; a >= startAngle; a -= step) {
            points.push({
                x: cx + innerRadius * Math.cos(a),
                y: cy + innerRadius * Math.sin(a),
            });
        }

        g.beginPath();
        g.moveTo(points[0].x, points[0].y);
        for (let p of points) g.lineTo(p.x, p.y);
        g.closePath();
        g.fillPath();
    }

    g.setBlendMode(Phaser.BlendModes.ADD);
    g.generateTexture(key, texSize, texSize);
    g.destroy();
}
