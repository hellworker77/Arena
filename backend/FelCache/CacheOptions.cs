namespace FelCache;

public record CacheOptions
{
    public string? BaseUrl { get; set; } = "http://localhost:5000";
    
    public UInt64 Ttl { get; set; } = 3600;
}