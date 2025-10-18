using Application.Contracts.Game;
using Application.Services.Game;
using Shared.Constants;

namespace Infrastructure.Services.Game;

public class CollisionService: ICollisionService
{
    public void Handle(ArenaState state)
    {
        var projectiles = state.Projectiles.ToList();
        var enemies = state.Enemies.ToList();

        foreach (var p in projectiles)
        {
            float prevX = p.X - p.Vx * ArenaGameConstants.PhysicsDelta;
            float prevY = p.Y - p.Vy * ArenaGameConstants.PhysicsDelta;

            bool hit = false;

            foreach (var e in enemies)
            {
                if (CheckSegmentCircleHit(prevX, prevY, p.X, p.Y, e.X, e.Y, ArenaGameConstants.EnemyRadius))
                {
                    e.Health -= p.Damage;
                    hit = true;

                    if (e.Health <= 0)
                        state.Enemies.Remove(e);

                    break;
                }
            }

            if (hit)
                state.Projectiles.Remove(p);
        }
    }

    private bool CheckSegmentCircleHit(float x1,
        float y1,
        float x2,
        float y2,
        float cx,
        float cy,
        float radius)
    {
        float dx = x2 - x1;
        float dy = y2 - y1;

        float fx = cx - x1;
        float fy = cy - y1;

        float t = (fx * dx + fy * dy) / (dx * dx + dy * dy);
        t = Math.Clamp(t, 0f, 1f);

        float closestX = x1 + dx * t;
        float closestY = y1 + dy * t;

        float distSq = (closestX - cx) * (closestX - cx) + (closestY - cy) * (closestY - cy);
        return distSq <= radius * radius;
    }
}