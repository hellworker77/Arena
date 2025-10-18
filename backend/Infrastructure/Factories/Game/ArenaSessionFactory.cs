using Application.Contracts.Game;
using Application.Contracts.Game.Abstract;
using Application.EventBus.Notifiers;
using Application.Factories.Game;
using Application.Services.Game;
using Infrastructure.Contracts.Models;

namespace Infrastructure.Factories.Game;

public class ArenaSessionFactory(IArenaWorldService worldService,
    IStateNotifier<ArenaState> stateNotifier): IArenaSessionFactory
{
    public IArenaSession Create(string connectionId)
        => new ArenaSession(connectionId, worldService, stateNotifier);
}