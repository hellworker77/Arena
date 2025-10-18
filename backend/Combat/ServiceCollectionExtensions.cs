using Application.Contracts.Game;
using Application.EventBus;
using Application.EventBus.Notifiers;
using Application.Fabrics;
using Application.Factories.Game;
using Application.Repositories.Game;
using Application.Services.Game;
using Infrastructure.EventBus;
using Infrastructure.EventBus.Notifiers;
using Infrastructure.Fabrics.Arena;
using Infrastructure.Factories.Game;
using Infrastructure.Hubs;
using Infrastructure.Services.Game;
using Persistence.Contexts;
using Persistence.Repositories.Game;

namespace Combat;

public static class ServiceCollectionExtensions
{
    public static IServiceCollection AddServices(this IServiceCollection services)
        => services
            .AddScoped<IPlayerService, PlayerService>()
            .AddScoped<IEnemyService, EnemyService>()
            .AddScoped<IProjectileService, ProjectileService>()
            .AddScoped<ICollisionService, CollisionService>()
            .AddScoped<IArenaWorldService, ArenaWorldService>()
            .AddScoped<IArenaService, ArenaService>();
    
    public static IServiceCollection AddWebRtc(this IServiceCollection services)
        => services
            .AddScoped<ArenaWebRtcService>();
    
    public static IServiceCollection AddRepositories(this IServiceCollection services)
        => services
            .AddScoped<IArenaRepository, ArenaRepository>();
    
    public static IServiceCollection AddEventBus(this IServiceCollection services)
        => services
            .AddSingleton<IStateNotifier<ArenaState>, StateNotifier<ArenaState>>()
            .AddSingleton<IEventBus, InMemoryArenaEventBus>();
    
    public static IServiceCollection AddFactories(this IServiceCollection services)
        => services
            .AddScoped<IArenaMessageHandlerFactory, ArenaMessageHandlerFactory>()
            .AddScoped<IArenaSessionFactory, ArenaSessionFactory>();
    
    public static IServiceCollection AddFabrics(this IServiceCollection services)
        => services
            .AddScoped<IArenaClientMessageHandler, GetSnapshotHandler>()
            .AddScoped<IArenaClientMessageHandler, MoveHandler>();
    
    public static IServiceCollection AddStores(this IServiceCollection services)
        => services
            .AddSingleton<InMemoryGameStore>();
}