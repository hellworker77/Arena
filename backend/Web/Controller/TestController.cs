using FelCache;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;

namespace Web.Controller;

[ApiController]
[Route("api/[controller]")]
public class TestController() : ControllerBase
{
    [Authorize(Roles = "Admin")]
    [HttpGet]
    public int Test()
    {
        return 5;
    }

    [HttpGet("key")]
    public string GetFromCache(string key)
    {
        return "Hexadecimal";
    }
}