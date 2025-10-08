interface Entity {
    health: number,
    defense: number,
    damage: { min: number, max: number },
    speed: number,
    attackSpeed: number,
    regeneration: number,

    isMelee: boolean,
    attackRange: number,
}

const playerConfig: Entity = {
    health: 626,
    defense: 4,
    damage: { min: 54, max: 60 },
    speed: 100,
    attackSpeed: 1000,
    regeneration: 2.4,

    isMelee: false,
    attackRange: 800,
}

const enemyConfig: Entity = {
    health: 100,
    defense: 1,
    damage: { min: 5, max: 12 },
    speed: 50,
    attackSpeed: 1000,
    regeneration: 0,

    isMelee: true,
    attackRange: 50,
}

export const gameOptions = {
    gameSize: {
        width               : 800,
        height              : 600
    },
    gameBackgroundColor     : 0x222,

    player                  : playerConfig,

    enemy                   : enemyConfig,

    enemySpawnRate          : 1000,

    projectile              : {
        velocity            : 425,
    }
}