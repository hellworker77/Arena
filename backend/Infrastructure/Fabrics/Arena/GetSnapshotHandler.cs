using System.Text.Json;
using Application.EventBus;
using Application.EventBus.Events;
using Application.Fabrics;
using Application.Services.Game;

namespace Infrastructure.Fabrics.Arena;

public class GetSnapshotHandler(IArenaService arenaService,
    IEventBus eventBus) : IArenaClientMessageHandler
{
    public string Type => "getSnapshot";
    
    public async Task HandleAsync(string connectionId, JsonElement payload)
    {
        var state = await arenaService.GetStateAsync(connectionId);
        eventBus.Publish(new SnapshotReadyEvent(connectionId, state));
    }
}