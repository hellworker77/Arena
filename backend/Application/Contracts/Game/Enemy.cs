using Application.Contracts.Game.Abstract;
using Shared.Constants;

namespace Application.Contracts.Game;

public class Enemy: AliveGameEntity
{
    public override int MaxHealth { get; set; } = ArenaGameConstants.EnemyHealth;
    
    public override int Health { get; set; } = ArenaGameConstants.EnemyHealth;
}