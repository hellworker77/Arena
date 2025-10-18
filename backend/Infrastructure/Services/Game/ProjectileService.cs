using Application.Contracts.Game;
using Application.Services.Game;
using Shared.Constants;

namespace Infrastructure.Services.Game;

public class ProjectileService: IProjectileService
{
    public void MoveProjectiles(List<Projectile> projectiles,
        float deltaTime)
    {
        foreach (var p in projectiles.ToArray())
        {
            p.X += p.Vx * deltaTime;
            p.Y += p.Vy * deltaTime;
            p.Lifespan -= deltaTime;

            if (p.Lifespan <= 0)
                projectiles.Remove(p);
        }
    }

    public void ShootAtNearestEnemy(ArenaState state)
    {
        if (state.Enemies.Count == 0) return;

        Enemy? closest = null;
        float minDist = float.MaxValue;

        foreach (var e in state.Enemies)
        {
            float dx = e.X - state.Player.X;
            float dy = e.Y - state.Player.Y;
            float dist = dx * dx + dy * dy;
            if (dist < minDist)
            {
                minDist = dist;
                closest = e;
            }
        }

        if (closest == null) return;

        float dxDir = closest.X - state.Player.X;
        float dyDir = closest.Y - state.Player.Y;
        float distDir = MathF.Sqrt(dxDir * dxDir + dyDir * dyDir);

        var projectile = new Projectile
        {
            Id = Guid.NewGuid(),
            X = state.Player.X,
            Y = state.Player.Y,
            Vx = dxDir / distDir * ArenaGameConstants.ProjectileSpeed,
            Vy = dyDir / distDir * ArenaGameConstants.ProjectileSpeed 
        };

        state.Projectiles.Add(projectile);
    }
}