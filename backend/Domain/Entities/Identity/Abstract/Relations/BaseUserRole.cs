using Domain.Entities.Abstract;

namespace Domain.Entities.Identity.Abstract.Relations;

public abstract class BaseUserRole<TUser, TRole> : BaseEntity
{
    public Guid UserId { get; set; }
    public virtual TUser User { get; set; }
    
    public Guid RoleId { get; set; }
    public virtual TRole Role { get; set; }
}