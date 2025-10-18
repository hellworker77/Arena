using Application.Contracts.Game;

namespace Application.Services.Game;

public interface IProjectileService
{
    void MoveProjectiles(List<Projectile> projectiles,
        float deltaTime);

    void ShootAtNearestEnemy(ArenaState state);
}