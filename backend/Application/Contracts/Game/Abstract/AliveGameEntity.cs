namespace Application.Contracts.Game.Abstract;

public abstract class AliveGameEntity: GameEntity
{
    public virtual int MaxHealth { get; set; }
    
    public virtual int Health { get; set; }
    
    public float TargetX { get; set; }
    
    public float TargetY { get; set; }
}