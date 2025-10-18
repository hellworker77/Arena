using Application.Contracts.Game;

namespace Application.Services.Game;

public interface IEnemyService
{
    void MoveEnemies(List<Enemy> enemies,
        Player player,
        float deltaTime);

    void SpawnEnemy(ArenaState state);
}