using Application.Contracts.Game;
using Application.EventBus;
using Application.EventBus.Events;
using Infrastructure.Hubs;
using Microsoft.AspNetCore.SignalR;

namespace Infrastructure.Contracts.Observers;

public class ArenaStateObserver(string connectionId,
    IEventBus eventBus): IObserver<ArenaState>
{
    public void OnCompleted()
    {
        Console.WriteLine($"[Observer] Session {connectionId} completed");
    }

    public void OnError(Exception error)
    {
        Console.WriteLine($"[Observer] Error for {connectionId}: {error}");
    }

    public void OnNext(ArenaState state)
    {
        eventBus.Publish(new SnapshotReadyEvent(connectionId, state));
    }
}