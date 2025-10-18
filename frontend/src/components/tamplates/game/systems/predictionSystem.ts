export class PredictionSystem {
    static lerp(a: number, b: number, t: number): number {
        return a + (b - a) * t;
    }

    // Экспоненциальное сглаживание с учётом dt (в секундах).
    // smoothingTime — время стабилизации в секундах (меньше = жёстче/быстрее).
    // deadzone — если расстояние меньше, просто ставим цель (чтобы убрать мелкие дерганья).
    static smoothPosition(current: number, target: number, dt: number, smoothingTime = 0.1, deadzone = 0.5): number {
        const diff = target - current;
        if (Math.abs(diff) <= deadzone) return target;
        // Альфа по формуле экспоненциального затухания: 1 - exp(-dt / tau)
        const alpha = 1 - Math.exp(-Math.max(dt, 0) / Math.max(smoothingTime, 1e-6));
        return current + diff * alpha;
    }

    // Экстраполяция позиции по скорости (msSinceServer в миллисекундах)
    static predictPosition(x: number, y: number, vx: number, vy: number, msSinceServer: number): { x: number; y: number } {
        const s = msSinceServer / 1000;
        return { x: x + vx * s, y: y + vy * s };
    }

    static distance(x1: number, y1: number, x2: number, y2: number): number {
        const dx = x2 - x1;
        const dy = y2 - y1;
        return Math.sqrt(dx * dx + dy * dy);
    }

    static interpolateColor(healthRatio: number): number {
        const color = Phaser.Display.Color.Interpolate.ColorWithColor(
            new Phaser.Display.Color(0, 0, 0),
            new Phaser.Display.Color(255, 255, 255),
            1,
            Phaser.Math.Clamp(healthRatio, 0, 1)
        );
        return Phaser.Display.Color.GetColor(color.r, color.g, color.b);
    }
}