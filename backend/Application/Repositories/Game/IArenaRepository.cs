using Application.Contracts.Game.Abstract;

namespace Application.Repositories.Game;

public interface IArenaRepository
{
    public IArenaSession GetOrCreateSession(string connectionId);
    
    bool CloseSession(string connectionId);
}