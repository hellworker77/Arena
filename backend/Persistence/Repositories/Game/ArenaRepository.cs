using System.Collections.Concurrent;
using Application.Contracts.Game;
using Application.Contracts.Game.Abstract;
using Application.Factories.Game;
using Application.Repositories.Game;
using Persistence.Contexts;

namespace Persistence.Repositories.Game;

public class ArenaRepository(InMemoryGameStore store,
    IArenaSessionFactory arenaSessionFactory): IArenaRepository
{
    public IArenaSession GetOrCreateSession(string connectionId)
        => store.Sessions.GetOrAdd(connectionId, _ => arenaSessionFactory.Create(connectionId));

    public bool CloseSession(string connectionId)
    {
        if (store.Sessions.TryRemove(connectionId, out var session))
        {
            session.Stop();
            return true;
        }

        return false;
    }
}