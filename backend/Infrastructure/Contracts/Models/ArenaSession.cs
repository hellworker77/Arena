using Application.Contracts.Game;
using Application.Contracts.Game.Abstract;
using Application.EventBus.Notifiers;
using Application.Services.Game;
using Shared.Constants;

namespace Infrastructure.Contracts.Models;

public class ArenaSession : IArenaSession
{
    private readonly IArenaWorldService _world;
    private readonly IStateNotifier<ArenaState> _notifier;
    private readonly CancellationTokenSource _cts = new();

    private readonly string _connectionId;
    
    public ArenaState State { get; }

    public ArenaSession(
        string connectionId,
        IArenaWorldService world,
        IStateNotifier<ArenaState> notifier)
    {
        _connectionId = connectionId;
        _world = world;
        _notifier = notifier;
        
        State = new ArenaState
        {
            Player = new Player { Id = Guid.NewGuid(), X = 100, Y = 100 },
            Enemies = new List<Enemy>(),
            Projectiles = new List<Projectile>()
        };
        
        Task.Run(UpdateLoopAsync);
    }

    public void Stop() => _cts.Cancel();

    public IDisposable Subscribe(IObserver<ArenaState> observer)
        => _notifier.Subscribe(observer, State);

    public void SetPlayerTarget(float x, float y)
        => _world.PlayerService.SetTargetPosition(State.Player, x, y);

    private async Task UpdateLoopAsync()
    {
        var token = _cts.Token;
        float sendTimer = 0f;
        
        Console.WriteLine($"[ArenaSession:{_connectionId}] started loop");

        while (!token.IsCancellationRequested)
        {
            try
            {
                _world.Update(State, ArenaGameConstants.PhysicsDelta);

                sendTimer += ArenaGameConstants.PhysicsDelta;
                if (sendTimer >= ArenaGameConstants.NotifyDelta)
                {
                    _notifier.Notify(State);
                    sendTimer = 0f;
                }

                await Task.Delay((int)(ArenaGameConstants.PhysicsDelta * 1000), token);
            }
            catch (TaskCanceledException)
            {
                break;
            }
            catch (Exception ex)
            {
                Console.WriteLine($"[ArenaSession:{_connectionId}] Error: {ex}");
                _notifier.NotifyError(ex);
                break;
            }
        }

        _notifier.NotifyCompleted();
        Console.WriteLine($"[ArenaSession:{_connectionId}] stopped loop");
    }
}