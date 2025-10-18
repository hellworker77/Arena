using Application.Fabrics;

namespace Application.Factories.Game;

public interface IArenaMessageHandlerFactory
{
    IArenaClientMessageHandler? CreateHandler(string type);
}