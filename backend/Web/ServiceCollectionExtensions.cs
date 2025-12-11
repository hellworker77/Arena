using ApiGateway.Migrations;
using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.EntityFrameworkCore;
using Persistence.Contexts;

namespace Web;

internal static class ServiceCollectionExtensions
{
    extension(IServiceCollection services)
    {
        internal IServiceCollection AddFelAuthentication(IConfiguration configuration)
        {
            var authority = configuration["Identity:Authority"];
            var audience = configuration["Identity:Audience"];
        
            services
                .AddAuthentication(JwtBearerDefaults.AuthenticationScheme)
                .AddJwtBearer(options =>
                {
                    options.Authority = authority; 
                    options.Audience = audience;
                    options.RequireHttpsMetadata = false;
                });
        
            return services;
        }

        internal void AddDbContexts(IConfiguration configuration)
        {
            var connectionString = configuration.GetConnectionString("masterDb")
                                   ?? throw new InvalidOperationException("Connection string 'masterDb' not found.");

            services.AddDbContext<ApplicationDbContext>(options =>
            {
                options.UseNpgsql(connectionString, 
                    builder => builder.MigrationsAssembly(typeof(Init).Assembly.FullName));
            
                options.UseLazyLoadingProxies();
            });
        }
    }
}