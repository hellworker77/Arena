using Application.DatabaseBootstrappers;
using Application.DatabaseFakers;
using Microsoft.EntityFrameworkCore;
using Persistence.Contexts;
using Persistence.DatabaseBootstrappers;
using Persistence.DatabaseFakers;

namespace ApiGateway;

internal static class ServiceCollectionExtensions
{
    internal static IServiceCollection AddDbContexts(this IServiceCollection services,
        IConfiguration configuration)
    {
        var connectionString = configuration.GetConnectionString("masterDb")
                               ?? throw new InvalidOperationException("Connection string 'masterDb' not found.");

        services.AddScoped<IApplicationDbContextBootstrapper, ApplicationDbContextBootstrapper>();
        services.AddScoped<IDbInitializer, DbInitializer>();
        
        return services.AddDbContext<ApplicationDbContext>(options =>
        {
            options.UseNpgsql(connectionString, builder => builder.MigrationsAssembly(typeof(Program).Assembly.FullName));
            
            options.UseLazyLoadingProxies();
        });
    }
}