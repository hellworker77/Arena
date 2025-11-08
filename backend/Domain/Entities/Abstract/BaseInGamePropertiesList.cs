namespace Domain.Entities.Abstract;

public abstract class BaseInGamePropertiesList
{
    #region Base Stats

    /// <summary>
    /// Defense value of the item
    /// </summary>
    public int Defense { get; set; }
    
    /// <summary>
    /// Minimum flat damage the item can deal
    /// </summary>
    public int MinimumDamage { get; set; }
    
    /// <summary>
    /// Maximum flat damage the item can deal
    /// </summary>
    public int MaximumDamage { get; set; }
    
    /// <summary>
    /// Base attack speed of the item
    /// </summary>
    public int AttackSpeed { get; set; }
    
    /// <summary>
    /// Level required to use the item
    /// </summary>
    public int RequiredLevel { get; set; }
    
    /// <summary>
    /// Defense pierce value of the item
    /// </summary>
    public int AttackRating { get; set; }

    #endregion

    #region Additional Attributes

    /// <summary>
    /// Additional strength provided by the item
    /// </summary>
    public int AdditionalStrength { get; set; }
    
    /// <summary>
    /// Additional agility provided by the item
    /// </summary>
    public int AdditionalAgility { get; set; }
    
    /// <summary>
    /// Additional intelligence provided by the item
    /// </summary>
    public int AdditionalIntelligence { get; set; }

    #endregion

    #region Additional Flat Bonuses

    /// <summary>
    /// Additional minimum damage the item can deal
    /// </summary>
    public int AdditionalMinimumDamage { get; set; }
    
    /// <summary>
    /// Additional maximum damage the item can deal
    /// </summary>
    public int AdditionalMaximumDamage { get; set; }
    
    /// <summary>
    /// Additional defense pierce value provided by the item
    /// </summary>
    public int AdditionalAttackRating { get; set; }

    #endregion

    #region Percent Bonuses

    /// <summary>
    /// Enhanced damage percentage provided by the item
    /// </summary>
    public float EnhancedDamage { get; set; }
    
    /// <summary>
    /// Enhanced defense percentage provided by the item
    /// </summary>
    public float EnhancedDefense { get; set; }
    
    /// <summary>
    /// Increased chance of blocking an attack provided by the item
    /// </summary>
    public float IncreaseChanceOfBlocking { get; set; }
    
    /// <summary>
    /// Life steal percentage provided by the item
    /// </summary>
    public float LifeSteal { get; set; }
    
    /// <summary>
    /// Mana steal percentage provided by the item
    /// </summary>
    public float ManaSteal { get; set; }
    
    /// <summary>
    /// Damage reduction percentage provided by the item
    /// </summary>
    public float DamageReduced { get; set; }

    #endregion

    #region Resistances

    /// <summary>
    /// Fire resistance percentage provided by the item
    /// </summary>
    public float FireResistance { get; set; }
    
    /// <summary>
    /// Cold resistance percentage provided by the item
    /// </summary>
    public float ColdResistance { get; set; }
    
    /// <summary>
    /// Lightning resistance percentage provided by the item
    /// </summary>
    public float LightningResistance { get; set; }
    
    /// <summary>
    /// Poison resistance percentage provided by the item
    /// </summary>
    public float PoisonResistance { get; set; }

    #endregion

    #region Absorbs

    /// <summary>
    /// Fire absorb percentage provided by the item
    /// </summary>
    public float FireAbsorb { get; set; }
    
    /// <summary>
    /// Cold absorb percentage provided by the item
    /// </summary>
    public float ColdAbsorb { get; set; }
    
    /// <summary>
    /// Lightning absorb percentage provided by the item
    /// </summary>
    public float LightningAbsorb { get; set; }
    
    /// <summary>
    /// Poison absorb percentage provided by the item
    /// </summary>
    public float PoisonAbsorb { get; set; }

    #endregion

    #region Defensive Effects

    /// <summary>
    /// Damage reflected to attacker percentage provided by the item
    /// </summary>
    public float AttackerTakesDamage { get; set; }

    /// <summary>
    /// Health replenishment rate provided by the item
    /// </summary>
    public float ReplenishLife { get; set; }

    #endregion

    #region Speed and Mobility

    /// <summary>
    /// Increased attack speed percentage provided by the item
    /// </summary>
    public float IncreaseAttackSpeed { get; set; }
    
    /// <summary>
    /// Increased movement speed percentage provided by the item
    /// </summary>
    public float FasterRunWalk { get; set; }
    
    /// <summary>
    /// Increased casting speed percentage provided by the item
    /// </summary>
    public float FasterCastRate { get; set; }

    #endregion

    #region Offensive Bonuses

    /// <summary>
    /// Bonus to attack rating provided by the item
    /// </summary>
    public float BonusToAttackRating { get; set; }
    
    /// <summary>
    /// Chance to open wounds percentage provided by the item
    /// </summary>
    public float ChanceToOpenWounds { get; set; }
    
    /// <summary>
    /// Chance to deliver a crushing blow percentage provided by the item
    /// </summary>
    public float ChanceToCrushingBlow { get; set; }

    /// <summary>
    /// Chance to deliver a deadly strike percentage provided by the item
    /// </summary>
    public float DeadlyStrike { get; set; }
    
    /// <summary>
    /// Critical strike chance percentage provided by the item
    /// </summary>
    public float CriticalStrike { get; set; }
    
    /// <summary>
    /// Chance to pierce percentage provided by the item
    /// </summary>
    public float PierceChance { get; set; }

    #endregion

    #region Elemental Skill Damage

    /// <summary>
    /// Fire skill damage bonus percentage provided by the item
    /// </summary>
    public float FireSkillDamage { get; set; }
    
    /// <summary>
    /// Cold skill damage bonus percentage provided by the item
    /// </summary>
    public float ColdSkillDamage { get; set; }
    
    /// <summary>
    /// Lightning skill damage bonus percentage provided by the item
    /// </summary>
    public float LightningSkillDamage { get; set; }
    
    /// <summary>
    /// Poison skill damage bonus percentage provided by the item
    /// </summary>
    public float PoisonSkillDamage { get; set; }

    #endregion
}