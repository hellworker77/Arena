using Domain.Entities.Abstract;

namespace Domain.Entities.Relations;

public class ApplicationUserRole: BaseEntity
{
    public Guid ApplicationUserId { get; set; }
    
    public virtual ApplicationUser ApplicationUser { get; set; }
    
    public Guid ApplicationRoleId { get; set; }
    
    public virtual ApplicationRole ApplicationRole { get; set; }
}