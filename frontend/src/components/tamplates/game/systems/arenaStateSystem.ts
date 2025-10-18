// typescript
import type { ArenaState } from "../../../../@types/game/arenaState";
import {arenaWebRTCService} from "../../../../ws/arenaWebRTCService.ts";

type MaybeServerState = ArenaState & {
    serverTimestamp?: number;
    ts?: number;
};

type WithVel = { vx?: number; vy?: number };

export class ArenaStateSystem {
    private buffer: {state: ArenaState; timestamp: number}[] = [];
    private delay: number = 0.12; // 120 ms для более мягкой интерполяции
    private lastState: ArenaState | null = null;
    private serverTimeOffset: number | null = null; // arrival - serverTs
    private maxBufferMs = 2000; // хранить до 2s данных
    private maxExtrapolateMs = 200; // не экстраполировать дальше 200ms

    public onStateUpdated?: (state: ArenaState) => void;

    updateFromServer (state: ArenaState): void {
        const arrival = performance.now();
        const serverState = state as MaybeServerState;
        const serverTs = serverState.serverTimestamp ?? serverState.ts ?? null;

        if (serverTs != null && this.serverTimeOffset == null) {
            this.serverTimeOffset = arrival - serverTs;
        }

        const entryTs = serverTs != null && this.serverTimeOffset != null
            ? serverTs + this.serverTimeOffset
            : arrival;

        this.buffer.push({state, timestamp: entryTs});
        this.lastState = state;

        const cutoff = arrival - this.delay * 1000 - this.maxBufferMs;
        this.buffer = this.buffer.filter(entry => entry.timestamp >= cutoff);
        if (this.buffer.length > 50) this.buffer.splice(0, this.buffer.length - 50);

        this.onStateUpdated?.(state);
    }

    private mergeAndInterpolateProjectiles(
        older: ArenaState["projectiles"],
        newer: ArenaState["projectiles"],
        t?: number,
        extrapolateS?: number
    ): ArenaState["projectiles"] {
        const map = new Map<string, ArenaState["projectiles"][number]>();

        older.forEach(p => map.set(p.id, { ...p }));

        newer.forEach(np => {
            const existing = map.get(np.id);
            if (existing) {
                if (typeof t === "number") {
                    existing.x = Phaser.Math.Linear(existing.x, np.x, t);
                    existing.y = Phaser.Math.Linear(existing.y, np.y, t);
                } else if (typeof extrapolateS === "number") {
                    const ev = existing as unknown as WithVel;
                    existing.x = existing.x + (ev.vx ?? 0) * extrapolateS;
                    existing.y = existing.y + (ev.vy ?? 0) * extrapolateS;
                }

                existing.vx = np.vx ?? existing.vx;
                existing.vy = np.vy ?? existing.vy;
                map.set(existing.id, existing);
            } else {
                map.set(np.id, { ...np });
            }
        });

        return Array.from(map.values());
    }

    private mergeAndInterpolateEnemies(
        older: ArenaState["enemies"],
        newer: ArenaState["enemies"],
        t: number
    ): ArenaState["enemies"] {
        const map = new Map<string, ArenaState["enemies"][number]>();
        older.forEach(e => map.set(e.id, { ...e }));
        newer.forEach(ne => {
            const ex = map.get(ne.id);
            if (ex) {
                ex.x = Phaser.Math.Linear(ex.x, ne.x, t);
                ex.y = Phaser.Math.Linear(ex.y, ne.y, t);
                ex.health = ne.health;
                ex.maxHealth = ne.maxHealth;
                map.set(ex.id, ex);
            } else {
                map.set(ne.id, { ...ne });
            }
        });
        return Array.from(map.values());
    }

    getInterpolatedState(): ArenaState | null {
        const now = performance.now() - this.delay * 1000;
        if (this.buffer.length === 0) {
            if (!this.lastState) return null;
            return JSON.parse(JSON.stringify(this.lastState));
        }

        if (this.buffer.length === 1) {
            const only = this.buffer[0];
            const dtMs = Math.min(Math.max(0, now - only.timestamp), this.maxExtrapolateMs);
            const s = dtMs / 1000;
            const clone: ArenaState = JSON.parse(JSON.stringify(only.state));

            const pv = clone.player as unknown as WithVel;
            if (pv.vx != null || pv.vy != null) {
                clone.player.x += (pv.vx ?? 0) * s;
                clone.player.y += (pv.vy ?? 0) * s;
            }

            clone.projectiles = clone.projectiles.map((p: ArenaState["projectiles"][number]) => ({ ...p }));
            clone.projectiles = this.mergeAndInterpolateProjectiles(clone.projectiles, [], undefined, s);

            return clone;
        }

        let older = this.buffer[0], newer = this.buffer[this.buffer.length - 1];
        for (let i = 0; i < this.buffer.length - 1; i++) {
            const s0 = this.buffer[i], s1 = this.buffer[i + 1];
            if (now >= s0.timestamp && now <= s1.timestamp) {
                older = s0; newer = s1; break;
            }
        }

        if (now > newer.timestamp) {
            const dtExtra = Math.min(now - newer.timestamp, this.maxExtrapolateMs);
            const s = dtExtra / 1000;
            const base: ArenaState = JSON.parse(JSON.stringify(newer.state));

            const pv = base.player as unknown as WithVel;
            if (pv.vx != null || pv.vy != null) {
                base.player.x += (pv.vx ?? 0) * s;
                base.player.y += (pv.vy ?? 0) * s;
            }

            base.projectiles = base.projectiles.map((p: ArenaState["projectiles"][number]) => ({ ...p }));
            base.projectiles = this.mergeAndInterpolateProjectiles(base.projectiles, [], undefined, s);

            return base;
        }

        const denom = (newer.timestamp - older.timestamp);
        const rawT = denom > 0 ? (now - older.timestamp) / denom : 0;
        const t = Phaser.Math.Clamp(rawT, 0, 1);

        const player = {
            ...older.state.player,
            x: Phaser.Math.Linear(older.state.player.x, newer.state.player.x, t),
            y: Phaser.Math.Linear(older.state.player.y, newer.state.player.y, t)
        };

        const enemies = this.mergeAndInterpolateEnemies(older.state.enemies, newer.state.enemies, t);

        const projectiles = this.mergeAndInterpolateProjectiles(older.state.projectiles, newer.state.projectiles, t);

        return {
            player,
            enemies,
            projectiles
        } as ArenaState;
    }

    movePlayerTo(x: number, y: number) {
        arenaWebRTCService.moveTo(x, y);
    }
}
