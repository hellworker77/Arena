using Application.Contracts.Game;
using Application.Services.Game;
using Shared.Constants;

namespace Infrastructure.Services.Game;

public class PlayerService: IPlayerService
{
    public void Move(Player player, float deltaTime)
    {
        float dx = player.TargetX - player.X;
        float dy = player.TargetY - player.Y;
        float dist = MathF.Sqrt(dx * dx + dy * dy);
        if (dist < 0.01f) return;

        float dirX = dx / dist;
        float dirY = dy / dist;
        float move = ArenaGameConstants.PlayerMovementSpeed * deltaTime;

        if (move >= dist)
        {
            player.X = player.TargetX;
            player.Y = player.TargetY;
        }
        else
        {
            player.X += dirX * move;
            player.Y += dirY * move;
        }
    }

    public void SetTargetPosition(Player player,
        float x,
        float y)
    {
        player.TargetX = x;
        player.TargetY = y;
    }
}