using Application.Contracts.Data.Identity;
using Application.Services.Identity;
using Microsoft.AspNetCore.Mvc;

namespace Identity.Controllers;

[ApiController]
[Route("api/[controller]")]
public class TokenController(ITokenService tokenService) : ControllerBase
{
    /// <summary>
    /// Creates a new JWT token for the authenticated user.
    /// </summary>
    [HttpPost]
    public async Task<Jwt> CreateToken([FromBody] BaseLoginDto login,
        CancellationToken ct)
        => await tokenService.CreateJwtAsync(login, ct);

    /// <summary>
    /// Refreshes an existing JWT token using a valid refresh token.
    /// </summary>
    [HttpPut("refresh")]
    public async Task<Jwt> RefreshToken([FromBody] string refreshToken,
        CancellationToken ct)
        => await tokenService.RefreshJwtAsync(refreshToken, ct);

    /// <summary>
    /// Revokes a JWT token, making it invalid for future use.
    /// </summary>
    [HttpDelete("revoke")]
    public async Task RevokeToken([FromBody] string refreshToken,
        CancellationToken ct)
        => await tokenService.RevokeRefreshTokenAsync(refreshToken, ct);
}