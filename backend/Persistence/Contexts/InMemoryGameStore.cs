using System.Collections.Concurrent;
using Application.Contracts.Game;
using Application.Contracts.Game.Abstract;

namespace Persistence.Contexts;

public class InMemoryGameStore
{
    // Key: connectionId
    public ConcurrentDictionary<string, IArenaSession> Sessions { get; set; } = new();
}