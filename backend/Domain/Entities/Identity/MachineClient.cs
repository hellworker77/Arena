using Domain.Entities.Abstract;

namespace Domain.Entities.Identity;

public class MachineClient: EventInfEntity
{
    /// <summary>
    /// Client identifier not FK
    /// </summary>
    public string ClientId { get; set; } = String.Empty;
    
    /// <summary>
    /// Client secret hash for security purposes
    /// </summary>
    public string ClientSecretHash { get; set; } = String.Empty;
    
    /// <summary>
    /// Description of the machine client
    /// </summary>
    public string? Description { get; set; }
    
    /// <summary>
    /// Tokens issued to this machine client
    /// </summary>
    public virtual IEnumerable<JwtToken>? Tokens { get; set; }
}