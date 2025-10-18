using Application.Contracts.Game.Abstract;
using Shared.Constants;

namespace Application.Contracts.Game;

public class Player: AliveGameEntity
{
    public override int MaxHealth { get; set; } = ArenaGameConstants.PlayerHealth;
    
    public override int Health { get; set; } = ArenaGameConstants.PlayerHealth;
}