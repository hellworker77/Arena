using Application.Contracts.Data.Identity;
using Application.Repositories.Identity;
using Application.Services.Identity;
using Infrastructure.Services.Identity;
using Microsoft.EntityFrameworkCore;
using Persistence.Contexts;
using Persistence.Repositories.Identity;

namespace Identity;

internal static class ServiceCollectionExtensions
{
    extension(IServiceCollection services)
    {
        internal IServiceCollection AddRsaPem(string basePath = "keys")
        {
            var privateKeyPath = Path.Combine(basePath, "private.pem");
            var publicKeyPath = Path.Combine(basePath, "public.pem");

            if (!File.Exists(privateKeyPath) || !File.Exists(publicKeyPath))
                throw new FileNotFoundException($"RSA key files not found in path: {Path.GetFullPath(basePath)}");

            var privatePem = File.ReadAllText(privateKeyPath);
            var publicPem = File.ReadAllText(publicKeyPath);

            var rsaPem = new RsaPem
            {
                privatePem = privatePem,
                publicPem = publicPem
            };

            services.AddSingleton(rsaPem);
            return services;
        }

        internal IServiceCollection AddRepositories()
            => services
                .AddScoped<IUserManager, UserManager>()
                .AddScoped<IRoleManager, RoleManager>()
                .AddScoped<IUserRoleManager, UserRoleManager>()
                .AddScoped<IMachineClientRepository, MachineClientRepository>()
                .AddScoped<IJwtTokenRepository, JwtTokenRepository>();

        internal IServiceCollection AddServices()
            => services
                .AddScoped<ITokenService, TokenService>();

        internal IServiceCollection AddDbContexts(IConfiguration configuration)
        {
            var connectionString = configuration.GetConnectionString("masterDb")
                                   ?? throw new InvalidOperationException("Connection string 'masterDb' not found.");

            return services.AddDbContext<ApplicationDbContext>(options =>
            {
                options.UseNpgsql(connectionString, builder => builder.MigrationsAssembly("Web"));
            
                options.UseLazyLoadingProxies();
            });
        }
    }
}