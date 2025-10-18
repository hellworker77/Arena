using Application.Contracts.Game;
using Application.Services.Game;
using Shared.Constants;

namespace Infrastructure.Services.Game;

public class ArenaWorldService(IEnemyService enemyService,
    IPlayerService playerService,
    ICollisionService collisionService,
    IProjectileService projectileService): IArenaWorldService
{
    private float _enemySpawnTimer;
    private float _playerShootTimer;
    
    public IPlayerService PlayerService => playerService;
    public IEnemyService EnemyService => enemyService;
    public IProjectileService ProjectileService => projectileService;
    public ICollisionService CollisionService => collisionService;

    public void Update(ArenaState state, float deltaTime)
    {
        playerService.Move(state.Player, deltaTime);
        enemyService.MoveEnemies(state.Enemies, state.Player, deltaTime);
        projectileService.MoveProjectiles(state.Projectiles, deltaTime);
        
        _enemySpawnTimer += deltaTime;
        _playerShootTimer += deltaTime;
        
        if (_enemySpawnTimer >= ArenaGameConstants.EnemySpawnInterval)
        {
            enemyService.SpawnEnemy(state);
            _enemySpawnTimer = 0f;
        }

        if (_playerShootTimer >= ArenaGameConstants.PlayerShootInterval)
        {
            projectileService.ShootAtNearestEnemy(state);
            _playerShootTimer = 0f;
        }

        collisionService.Handle(state);
    }
}