export function createCircleWithPlusCut(
    scene: Phaser.Scene,
    key: string,
    size: number = 128,
    color: number = 0x00ffff,
    glowRadius: number = 16,
    glowIntensity: number = 0.015,
    cutThickness: number = 20 // ширина выреза "+"
) {
    const texSize = size + glowRadius * 2;
    const cx = texSize / 2;
    const cy = texSize / 2;
    const radius = size / 2;

    // 1️⃣ — создаем HTMLCanvas
    const canvas = document.createElement('canvas');
    canvas.width = texSize;
    canvas.height = texSize;
    const ctx = canvas.getContext('2d')!;
    ctx.clearRect(0, 0, texSize, texSize);

    // 2️⃣ — рисуем свечение
    for (let r = glowRadius; r > 0; r--) {
        const alpha = glowIntensity * (r / glowRadius);
        ctx.beginPath();
        ctx.fillStyle = `rgba(${(color >> 16) & 255}, ${(color >> 8) & 255}, ${color & 255}, ${alpha})`;
        ctx.arc(cx, cy, radius + r, 0, Math.PI * 2);
        ctx.fill();
    }

    // 3️⃣ — рисуем основной круг
    ctx.beginPath();
    ctx.fillStyle = `rgba(${(color >> 16) & 255}, ${(color >> 8) & 255}, ${color & 255}, 1)`;
    ctx.arc(cx, cy, radius, 0, Math.PI * 2);
    ctx.fill();

    // 4️⃣ — вырезаем "плюс" (прозрачный)
    ctx.globalCompositeOperation = 'destination-out';
    const len = size * 1.2;
    const half = cutThickness / 2;

    // горизонтальный вырез
    ctx.fillRect(cx - len / 2, cy - half, len, cutThickness);
    // вертикальный вырез
    ctx.fillRect(cx - half, cy - len / 2, cutThickness, len);

    // вернём обычный режим рисования
    ctx.globalCompositeOperation = 'source-over';

    // 5️⃣ — сохраняем как текстуру в Phaser
    const texture = scene.textures.createCanvas(key, texSize, texSize);
    const tctx = texture.getContext();
    tctx.clearRect(0, 0, texSize, texSize);
    tctx.drawImage(canvas, 0, 0);
    texture.refresh();
}
