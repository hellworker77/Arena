using Application.DatabaseBootstrappers;
using Application.DatabaseFakers;
using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.EntityFrameworkCore;
using Persistence.Contexts;
using Persistence.DatabaseBootstrappers;
using Persistence.DatabaseFakers;

namespace Web;

internal static class ServiceCollectionExtensions
{
    internal static IServiceCollection AddFelAuthentication(this IServiceCollection services,
        IConfiguration configuration)
    {
        var authority = configuration["Identity:Authority"];
        var audience = configuration["Identity:Audience"];
        
        services
            .AddAuthentication(JwtBearerDefaults.AuthenticationScheme)
            .AddJwtBearer(options =>
            {
                options.Authority = authority; //now discovery endpoint will be picked up automatically
                options.Audience = audience;
                options.RequireHttpsMetadata = false; //on prod turns true
            });
            

        return services;
    }
    
    internal static IServiceCollection AddDbContexts(this IServiceCollection services,
        IConfiguration configuration)
    {
        var connectionString = configuration.GetConnectionString("masterDb")
                               ?? throw new InvalidOperationException("Connection string 'masterDb' not found.");

        services.AddScoped<IApplicationDbContextBootstrapper, ApplicationDbContextBootstrapper>();
        services.AddScoped<IDbInitializer, DbInitializer>();
        
        return services.AddDbContext<ApplicationDbContext>(options =>
        {
            options.UseNpgsql(connectionString, builder => builder.MigrationsAssembly(typeof(ApplicationDbContext).Assembly.FullName));
            
            options.UseLazyLoadingProxies();
        });
    }
}