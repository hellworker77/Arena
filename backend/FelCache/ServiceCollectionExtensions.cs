using Microsoft.Extensions.DependencyInjection;

namespace FelCache;

public static class ServiceCollectionExtensions
{
    public static IServiceCollection AddFelCache(this IServiceCollection services, Action<CacheOptions> configureOptions)
    {
        var options = new CacheOptions();
        configureOptions(options);
        services.AddSingleton(options);
        services.AddSingleton<FelCacheClient>();
        return services;
    }
}