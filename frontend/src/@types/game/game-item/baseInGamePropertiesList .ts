/** Идиоматичный TS-интерфейс с camelCase именами */
export interface BaseInGameProperties {
    // base stats
    defense: number;
    minimumDamage: number;
    maximumDamage: number;
    attackSpeed: number;
    requiredLevel: number;
    attackRating: number;

    // attributes
    additionalStrength: number;
    additionalAgility: number;
    additionalIntelligence: number;

    // flat bonuses
    additionalMinimumDamage: number;
    additionalMaximumDamage: number;
    additionalAttackRating: number;

    // percent bonuses
    enhancedDamage: number;
    enhancedDefense: number;
    increaseChanceOfBlocking: number;
    lifeSteal: number;
    manaSteal: number;
    damageReduced: number;

    // resistances
    fireResistance: number;
    coldResistance: number;
    lightningResistance: number;
    poisonResistance: number;

    // absorbs
    fireAbsorb: number;
    coldAbsorb: number;
    lightningAbsorb: number;
    poisonAbsorb: number;

    // defensive effects
    attackerTakesDamage: number;
    replenishLife: number;

    // speed / mobility
    increaseAttackSpeed: number;
    fasterRunWalk: number;
    fasterCastRate: number;

    // offensive bonuses
    bonusToAttackRating: number;
    chanceToOpenWounds: number;
    chanceToCrushingBlow: number;
    deadlyStrike: number;
    criticalStrike: number;
    pierceChance: number;

    // elemental skill damage
    fireSkillDamage: number;
    coldSkillDamage: number;
    lightningSkillDamage: number;
    poisonSkillDamage: number;
}