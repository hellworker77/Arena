using Application.Contracts.Game.Abstract;
using Shared.Constants;

namespace Application.Contracts.Game;

public class Projectile: GameEntity
{
    public int Damage { get; set; } = ArenaGameConstants.ProjectileDamage;
    
    public float Vx { get; set; }
    
    public float Vy { get; set; }
    
    public float Lifespan { get; set; } = ArenaGameConstants.ProjectileLifeSpan;
}