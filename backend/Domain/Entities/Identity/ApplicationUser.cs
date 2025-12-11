using Domain.Entities.Identity.Abstract;
using Domain.Entities.Identity.Relations;

namespace Domain.Entities.Identity;

public class ApplicationUser : BaseUser<ApplicationUserRole>
{
    public virtual List<SavedCharacter> Characters { get; set; }
}