using FelCache;
using Microsoft.AspNetCore.Mvc;

namespace Web.Controller;

[ApiController]
[Route("api/[controller]")]
public class TestController
{
    private readonly FelCacheClient _felCacheClient;
    
    public TestController(FelCacheClient client)
    {
        _felCacheClient = client;
    }
    
    
    [HttpGet]
    public int Test()
    {
        _felCacheClient.SetAsync("test-key", 25, null).Wait();
        var (found, value) = _felCacheClient.GetAsync<int>("test-key").Result;

        return value;
    }
}