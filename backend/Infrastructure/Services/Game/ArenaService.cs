using Application.Contracts.Game;
using Application.EventBus;
using Application.Repositories.Game;
using Application.Services.Game;
using Infrastructure.Contracts.Observers;

namespace Infrastructure.Services.Game;

public class ArenaService(IArenaRepository arenaRepository,
    IEventBus eventBus) : IArenaService
{
    public Task InitializeSessionAsync(string connectionId)
    {
        var session = arenaRepository.GetOrCreateSession(connectionId);
        
        var observer = new ArenaStateObserver(connectionId, eventBus);
        session.Subscribe(observer);
        
        return Task.CompletedTask;
    }

    public Task CloseSessionAsync(string connectionId)
    {
        if (arenaRepository.CloseSession(connectionId))
            Console.WriteLine($"Closed session {connectionId}");
        return Task.CompletedTask;
    }

    public Task<ArenaState> GetStateAsync(string connectionId)
    {
        var session = arenaRepository.GetOrCreateSession(connectionId);
        return Task.FromResult(session.State);
    }

    public Task MovePlayerAsync(string connectionId, float x, float y)
    {
        var session = arenaRepository.GetOrCreateSession(connectionId);
        session.SetPlayerTarget(x, y);
        return Task.CompletedTask;
    }
}