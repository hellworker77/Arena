using Application.Contracts.Data.Identity;
using Microsoft.EntityFrameworkCore;

namespace Application.Services.Identity;

public interface ITokenService
{
    Task<Jwks> GetJwksAsync();
    
    Task<Jwt> CreateJwtAsync(BaseLoginDto login,
        CancellationToken cancellationToken = default);
    
    Task<Jwt> RefreshJwtAsync(string refreshToken, 
        CancellationToken cancellationToken = default);
    
    Task RevokeRefreshTokenAsync(string refreshToken,
        CancellationToken cancellationToken = default);
    
    Task<Jwt> CreateMachineClientJwtAsync(ClientCredentialsDto clientCredentials,
        CancellationToken cancellationToken = default);
}