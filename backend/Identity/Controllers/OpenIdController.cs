using Application.Contracts.Data.Identity;
using Application.Services.Identity;
using Microsoft.AspNetCore.Mvc;

namespace Identity.Controllers;

[ApiController]
[Route(".well-known")]
public class OpenIdController(ITokenService tokenService): ControllerBase
{
    /// <summary>
    /// OpenID Connect configuration endpoint.
    /// </summary>
    [HttpGet("openid-configuration")]
    public IActionResult GetConfiguration()
    {
        var authority = $"{Request.Scheme}://{Request.Host}";
        
        return Ok(new
        {
            issuer = authority,
            jwks_uri = $"{authority}/.well-known/jwks.json",
            token_endpoint = $"{authority}/api/token",
            // Add other OpenID Connect configuration parameters as needed
        });
    }
    
    /// <summary>
    /// JWKS endpoint to retrieve the public keys for token validation.
    /// </summary>
    [HttpGet("jwks.json")]
    public async Task<Jwks> GetJwks()
        => await tokenService.GetJwksAsync();
}