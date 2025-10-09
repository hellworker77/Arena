using Domain.Entities.Relations;

namespace Domain.Entities.Abstract;

public abstract class BaseRole<TUserRole> : EventInfEntity
{
    public string Name { get; set; }

    public virtual IEnumerable<TUserRole> UserRoles { get; set; }
}