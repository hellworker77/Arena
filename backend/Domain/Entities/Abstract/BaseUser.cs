using Domain.Entities.Relations;

namespace Domain.Entities.Abstract;

public abstract class BaseUser<TUserRole>: EventInfEntity
{
    public string UserName { get; set; }
    
    public string Email { get; set; }
    
    public string PasswordHash { get; set; }
    
    public DateTime BannedUntil { get; set; }
    
    public virtual IEnumerable<TUserRole> UserRoles { get; set; }
    
    public virtual IEnumerable<JwtToken>? Tokens { get; set; }
}