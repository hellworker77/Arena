using Domain.Entities.Identity;

namespace Application.Repositories.Identity;

public interface IJwtTokenRepository
{
    Task<JwtToken?> GetByRefreshTokenHashAsync(string refreshTokenHash,
        CancellationToken cancellationToken = default);
    
    Task SaveJwtAsync(JwtToken token,
        CancellationToken cancellationToken = default);
    
    Task UpdateJwtAsync(JwtToken token,
        CancellationToken cancellationToken = default);
    
    /// <summary>
    /// Retains only the most recent tokens for a user, deleting older ones and
    /// returning the count of deleted tokens.
    /// </summary>
    Task<int> RetainOldTokensAsync(Guid userId,
        int keepCount = 5,
        CancellationToken cancellationToken = default);
    
    /// <summary>
    /// Retains only the most recent tokens for a machine client, deleting older ones and
    /// returning the count of deleted tokens.
    /// </summary>
    Task<int> RetainOldMachineClientTokensAsync(Guid machineClientId,
        int keepCount = 5,
        CancellationToken cancellationToken = default);
}