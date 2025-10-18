using System.Text.Json;
using Application.EventBus;
using Application.EventBus.Events;
using Application.Fabrics;
using Application.Services.Game;

namespace Infrastructure.Fabrics.Arena;

public class MoveHandler(IArenaService arenaService,
    IEventBus eventBus): IArenaClientMessageHandler
{
    public string Type => "move";
    
    public async Task HandleAsync(string connectionId, JsonElement payload)
    {
        var x = payload.GetProperty("x").GetSingle();
        var y = payload.GetProperty("y").GetSingle();
        await arenaService.MovePlayerAsync(connectionId, x, y);
        
        var state = await arenaService.GetStateAsync(connectionId);
        eventBus.Publish(new SnapshotReadyEvent(connectionId, state));
    }
}