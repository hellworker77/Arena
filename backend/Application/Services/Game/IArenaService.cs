using Application.Contracts.Game;

namespace Application.Services.Game;

public interface IArenaService
{
    Task InitializeSessionAsync(string connectionId);
    
    Task CloseSessionAsync(string connectionId);
    
    Task<ArenaState> GetStateAsync(string connectionId);
    
    Task MovePlayerAsync(string connectionId,
        float x,
        float y);
}