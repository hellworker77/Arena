using Microsoft.AspNetCore.Authentication.JwtBearer;
using Shared.Settings;

namespace QAService;

internal static class ServiceCollectionExtensions
{
    private static string ApiCorsPolicyName => "ApiCors";

    extension(IServiceCollection services)
    {
        internal IServiceCollection ConfigureCors(IConfiguration configuration)
        {
            var corsOptions = configuration
                .GetSection("Cors")
                .Get<ApiCorsOptions>();

            if (corsOptions == null)
            {
                throw new InvalidOperationException("Cors configuration section is missing or invalid.");
            }

            services.AddCors(options =>
            {
                options.AddPolicy(ApiCorsPolicyName, builder =>
                {
                    builder.WithOrigins(corsOptions.AllowedOrigins)
                        .AllowAnyHeader()
                        .AllowAnyMethod();
                });
            });

            return services;
        }

        internal IServiceCollection AddFelAuthentication(IConfiguration configuration)
        {
            var authority = configuration["Identity:Authority"];
            var audience = configuration["Identity:Audience"];

            services
                .AddAuthentication(JwtBearerDefaults.AuthenticationScheme)
                .AddJwtBearer("JwtBearer", options =>
                {
                    options.Authority = authority;
                    options.Audience = audience;
                    options.RequireHttpsMetadata = false;
                });

            services.AddAuthorization(options =>
            {
                options.AddPolicy("CanRunTests", 
                    policy => policy.RequireRole("Admin", "QA", "Tester"));
            });

            return services;
        }
    }

    internal static void UseApiCors(this WebApplication app)
        => app.UseCors(ApiCorsPolicyName);
}