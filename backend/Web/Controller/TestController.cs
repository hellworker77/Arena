using FelCache;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;

namespace Web.Controller;

[ApiController]
[Route("api/[controller]")]
public class TestController: ControllerBase
{
    private readonly FelCacheClient _felCacheClient;
    
    public TestController(FelCacheClient client)
    {
        _felCacheClient = client;
    }
    
    [Authorize(Roles = "Admin")]
    [HttpGet]
    public int Test()
    {
        return 5;
    }
}