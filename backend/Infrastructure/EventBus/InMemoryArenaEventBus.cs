using System.Collections.Concurrent;
using Application.EventBus;

namespace Infrastructure.EventBus;

public class InMemoryArenaEventBus: IEventBus
{
    private readonly ConcurrentDictionary<Type, List<Delegate>> _handlers = new();
    
    public void Publish<TEvent>(TEvent @event)
    {
        if (_handlers.TryGetValue(typeof(TEvent), out var handlers))
        {
            foreach (var handler in handlers.Cast<Action<TEvent>>())
            {
                handler(@event);
            }
        }
    }

    public void Subscribe<TEvent>(Action<TEvent> handler)
    {
        var list = _handlers.GetOrAdd(typeof(TEvent), _ => new List<Delegate>());
        list.Add(handler);
    }
}