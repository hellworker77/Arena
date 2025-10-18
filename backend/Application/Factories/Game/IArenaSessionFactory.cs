using Application.Contracts.Game.Abstract;

namespace Application.Factories.Game;

public interface IArenaSessionFactory
{
    public IArenaSession Create(string connectionId);
}