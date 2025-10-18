namespace Application.EventBus.Notifiers;

public interface IStateNotifier<TState>
{
    IDisposable Subscribe(IObserver<TState> observer,
        TState initialState);
    
    void Notify(TState state);
    
    void NotifyError(Exception exception);
    
    void NotifyCompleted();
}