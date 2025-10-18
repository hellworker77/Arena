namespace Application.Contracts.Game;

public class ArenaState
{
    public Player Player { get; set; }
    
    public List<Enemy> Enemies { get; set; }
    
    public List<Projectile> Projectiles { get; set; }
}