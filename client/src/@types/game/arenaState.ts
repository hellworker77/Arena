export interface ArenaState {
    player: Player
    enemies: Enemy[]
    projectiles: Projectile[]
}

export interface Player {
    targetY: number;
    targetX: number;
    id: string
    x: number
    y: number
    health: number
    maxHealth: number
}

export interface Enemy {
    id: string
    x: number
    y: number
    health: number
    maxHealth: number
}

export interface Projectile {
    id: string
    x: number
    y: number
    damage: number
    vx: number
    vy: number
}