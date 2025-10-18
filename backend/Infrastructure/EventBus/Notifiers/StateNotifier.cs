using Application.EventBus.Notifiers;

namespace Infrastructure.EventBus.Notifiers;

public class StateNotifier<TState> : IStateNotifier<TState>
    where TState : class
{
    private readonly List<IObserver<TState>> _observers = new();
    
    public IDisposable Subscribe(IObserver<TState> observer, TState initialState)
    {
        if (!_observers.Contains(observer))
        {
            _observers.Add(observer);
            observer.OnNext(initialState);
        }

        return new Unsubscriber(_observers, observer);
    }

    public void Notify(TState state)
    {
        foreach (var obs in _observers.ToArray())
            obs.OnNext(state);
    }

    public void NotifyError(Exception ex)
    {
        foreach (var obs in _observers.ToArray())
            obs.OnError(ex);
    }

    public void NotifyCompleted()
    {
        foreach (var obs in _observers.ToArray())
            obs.OnCompleted();
        _observers.Clear();
    }
    
    private class Unsubscriber(List<IObserver<TState>> observers, 
        IObserver<TState> observer)
        : IDisposable
    {
        public void Dispose()
        {
            if (observers.Contains(observer))
                observers.Remove(observer);
        }
    }
}