using Application.Contracts.Game;
using Application.Services.Game;
using Shared.Constants;

namespace Infrastructure.Services.Game;

public class EnemyService: IEnemyService
{
    private readonly Random _random = new();
    
    public void MoveEnemies(List<Enemy> enemies,
        Player player,
        float deltaTime)
    {
        foreach (var enemy in enemies)
        {
            float dx = player.X - enemy.X;
            float dy = player.Y - enemy.Y;
            
            float dist = MathF.Sqrt(dx * dx + dy * dy);
            
            if(dist < ArenaGameConstants.EnemyDistanceEpsilon) return;
            
            enemy.X += dx / dist * ArenaGameConstants.EnemyMovementSpeed * deltaTime;
            enemy.Y += dy / dist * ArenaGameConstants.EnemyMovementSpeed * deltaTime;
        }
    }

    public void SpawnEnemy(ArenaState state)
    {
        var enemy = new Enemy
        {
            Id = Guid.NewGuid(),
            X = _random.Next(0, 400),
            Y = _random.Next(0, 400)
        };
        state.Enemies.Add(enemy);
    }
}