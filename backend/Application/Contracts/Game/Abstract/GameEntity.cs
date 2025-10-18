namespace Application.Contracts.Game.Abstract;

public abstract class GameEntity
{
    public Guid Id { get; set; }
    
    public float X { get; set; }
    
    public float Y { get; set; }
}