using Domain.Entities.Abstract;

namespace Domain.Entities.Identity.Abstract;

public abstract class BaseRole<TUserRole> : EventInfEntity
{
    public string Name { get; set; }

    public virtual IEnumerable<TUserRole> UserRoles { get; set; }
}