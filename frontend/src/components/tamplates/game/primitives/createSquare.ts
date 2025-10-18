export function createSquare(
    scene: Phaser.Scene,
    key: string,
    size: number = 32,
    color: number = 0x00ffff,
) {
    const texSize = size;
    const cx = texSize / 2;
    const cy = texSize / 2;

    const canvas = document.createElement('canvas');
    canvas.width = texSize;
    canvas.height = texSize;
    const ctx = canvas.getContext('2d')!;
    ctx.clearRect(0, 0, texSize, texSize);

    ctx.fillStyle = `rgba(${(color >> 16) & 255}, ${(color >> 8) & 255}, ${color & 255}, 1)`;
    ctx.fillRect(cx - size / 2, cy - size / 2, size, size);

    const texture = scene.textures.createCanvas(key, texSize, texSize);
    const tctx = texture!.getContext();
    tctx.clearRect(0, 0, texSize, texSize);
    tctx.drawImage(canvas, 0, 0);
    texture!.refresh();
}
