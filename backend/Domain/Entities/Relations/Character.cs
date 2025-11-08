using Domain.Entities.Abstract;

namespace Domain.Entities;

public partial class Character
{
    public virtual Inventory Inventory { get; set; } = null!;
}