import type {BaseInGameProperties} from "./baseInGamePropertiesList .ts";
import {Rarity} from "./rarity.ts";

export interface Item extends BaseInGameProperties {
    name: string;
    rarity: Rarity;
}

export const mightySwordOfAzeroth: Item = {
    name: "Mighty Sword of Azeroth",
    rarity: Rarity.EPIC,
    minimumDamage: 15,
    maximumDamage: 25,
    defense: 10,
    additionalStrength: 5,
    lifeSteal: 2,
    // остальные свойства 0 по умолчанию
    attackSpeed: 0,
    requiredLevel: 0,
    attackRating: 0,
    additionalAgility: 0,
    additionalIntelligence: 0,
    additionalMinimumDamage: 0,
    additionalMaximumDamage: 0,
    additionalAttackRating: 0,
    enhancedDamage: 0,
    enhancedDefense: 0,
    increaseChanceOfBlocking: 0,
    manaSteal: 0,
    damageReduced: 0,
    fireResistance: 0,
    coldResistance: 0,
    lightningResistance: 0,
    poisonResistance: 0,
    fireAbsorb: 0,
    coldAbsorb: 0,
    lightningAbsorb: 0,
    poisonAbsorb: 0,
    attackerTakesDamage: 0,
    replenishLife: 0,
    increaseAttackSpeed: 0,
    fasterRunWalk: 0,
    fasterCastRate: 0,
    bonusToAttackRating: 0,
    chanceToOpenWounds: 0,
    chanceToCrushingBlow: 0,
    deadlyStrike: 0,
    criticalStrike: 0,
    pierceChance: 0,
    fireSkillDamage: 0,
    coldSkillDamage: 0,
    lightningSkillDamage: 0,
    poisonSkillDamage: 0
};