namespace Domain.Entities;

public partial class Inventory
{
    public Guid CharacterId { get; set; }
    
    public virtual Character Character { get; set; } = null!;
}