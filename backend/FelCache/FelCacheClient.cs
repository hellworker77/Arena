using System.Text.Json;
using Cache;
using Grpc.Net.Client;

namespace FelCache;

public class FelCacheClient
{
    private readonly Cache.Cache.CacheClient _grpcClient;
    private readonly GrpcChannel _channel;
    private readonly CacheOptions _options;

    public FelCacheClient(CacheOptions options)
    {
        _channel = GrpcChannel.ForAddress(options.BaseUrl, new GrpcChannelOptions
        {
            HttpHandler = new HttpClientHandler()
        });

        _options = options;
        
        _grpcClient = new Cache.Cache.CacheClient(_channel);
    }
    
    public async Task<bool> SetAsync<T>(string key, T value, UInt64? ttl)
    {
        var reply = await _grpcClient.SetAsync(new SetRequest
        {
            Key = key,
            Value = JsonSerializer.Serialize(value),
            TtlMilliseconds = ttl ?? _options.Ttl
        });
        
        return reply.Saved;
    }
    
    public async Task<(bool found, T? value)> GetAsync<T>(string key)
    {
        var reply = await _grpcClient.GetAsync(new GetRequest
        {
            Key = key
        });

        return (reply.Found, reply.Found ? JsonSerializer.Deserialize<T>(reply.Value)! : default);
    }
    
    public async Task<bool> DeleteAsync(string key)
    {
        var reply = await _grpcClient.DeleteAsync(new DeleteRequest { Key = key });
        return reply.Removed;
    }
    
    public async Task<bool> UpdateAsync<T>(string key, T value, ulong? ttlMilliseconds)
    {
        var reply = await _grpcClient.UpdateAsync(new UpdateRequest
        {
            Key = key,
            Value = JsonSerializer.Serialize(value),
            TtlMilliseconds = ttlMilliseconds ?? _options.Ttl
        });
        return reply.Updated;
    }

    public async Task<bool> ClearAsync()
    {
        var reply = await _grpcClient.ClearAsync(new ClearRequest());
        return reply.Cleared;
    }
}