using System.Text.Json;

namespace Application.Fabrics;

public interface IArenaClientMessageHandler
{
    string Type { get; }

    Task HandleAsync(string connectionId, 
        JsonElement payload);
}