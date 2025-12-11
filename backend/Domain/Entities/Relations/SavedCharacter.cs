using Domain.Entities.Identity;

namespace Domain.Entities;

public partial class SavedCharacter
{
    public Guid ApplicationUserId { get; set; }

    public virtual ApplicationUser ApplicationUser { get; set; } = null!;
}