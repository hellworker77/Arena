using Application.Contracts.Game;

namespace Application.Services.Game;

public interface ICollisionService
{
    void Handle(ArenaState state);
}