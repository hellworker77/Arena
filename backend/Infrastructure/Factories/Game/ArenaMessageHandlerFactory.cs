using Application.Fabrics;
using Application.Factories.Game;

namespace Infrastructure.Factories.Game;

public class ArenaMessageHandlerFactory(IEnumerable<IArenaClientMessageHandler> handlers): IArenaMessageHandlerFactory
{
    public IArenaClientMessageHandler? CreateHandler(string type)
    {
        return handlers.FirstOrDefault(x => x.Type.ToLower() == type.ToLower());
    }
}