using Application.Contracts.Game;

namespace Application.Services.Game;

public interface IPlayerService
{
    void Move(Player player,
        float deltaTime);

    void SetTargetPosition(Player player,
        float x,
        float y);
}