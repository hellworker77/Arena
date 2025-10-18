namespace Application.Contracts.Game.Abstract;

public interface IArenaSession: IObservable<ArenaState>
{
    ArenaState State { get; }
    void SetPlayerTarget(float x, float y);
    void Stop();
}