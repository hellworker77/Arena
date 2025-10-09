using Microsoft.AspNetCore.Authentication.JwtBearer;

namespace Web;

internal static class ServiceCollectionExtensions
{
    internal static IServiceCollection AddFelAuthentication(this IServiceCollection services,
        IConfiguration configuration)
    {
        var authority = configuration["Identity:Authority"];
        
        services
            .AddAuthentication(JwtBearerDefaults.AuthenticationScheme)
            .AddJwtBearer(options =>
            {
                options.Authority = authority; //now discovery endpoint will be picked up automatically
                options.Audience = "Microservices";
                options.RequireHttpsMetadata = false; //on prod turns true
            });
            

        return services;
    }
}