using Application.Contracts.Game;

namespace Application.Services.Game;

public interface IArenaWorldService
{
    IPlayerService PlayerService { get; }
    
    IEnemyService EnemyService { get; }
    
    IProjectileService ProjectileService { get; }
    
    ICollisionService CollisionService { get; }
    
    void Update(ArenaState state, 
        float deltaTime);
}