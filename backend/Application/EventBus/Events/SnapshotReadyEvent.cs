using Application.Contracts.Game;

namespace Application.EventBus.Events;

public record SnapshotReadyEvent(string ClientId, ArenaState State);