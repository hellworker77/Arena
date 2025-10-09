using Microsoft.AspNetCore.Mvc;

namespace Identity.Controllers;

[ApiController]
[Route(".well-known")]
public class OpenIdController: ControllerBase
{
    [HttpGet]
    public IActionResult GetConfiguration()
    {
        var authority = $"{Request.Scheme}://{Request.Host}";
        
        return Ok(new
        {
            issuer = authority,
            jwks_uri = $"{authority}/api/token/.well-known/jwks.json",
            token_endpoint = $"{authority}/api/token",
            // Add other OpenID Connect configuration parameters as needed
        });
    }
}