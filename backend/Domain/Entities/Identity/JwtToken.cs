using Domain.Entities.Abstract;

namespace Domain.Entities.Identity;

public class JwtToken: EventInfEntity
{
    /// <summary>
    /// Hash of the access token for security purposes
    /// </summary>
    public string AccessTokenHash { get; set; }
    
    /// <summary>
    /// Access token expiration date
    /// </summary>
    public DateTime AccessTokenExpiresAt { get; set; }
    
    /// <summary>
    /// Indicates whether the token has been revoked
    /// </summary>
    public bool IsRevoked { get; set; }
    
    /// <summary>
    /// Refresh token hash for security purposes
    /// </summary>
    public string RefreshTokenHash { get; set; } = string.Empty;
    
    /// <summary>
    /// Refresh token expiration date
    /// </summary>
    public DateTime RefreshTokenExpiresAt { get; set; }
    
    /// <summary>
    /// Foreign key to the user
    /// </summary>
    public Guid? UserId { get; set; }
    
    /// <summary>
    /// Application user associated with this token
    /// </summary>
    public virtual ApplicationUser? User { get; set; }
    
    /// <summary>
    /// Foreign key to the machine client
    /// </summary>
    public Guid? MachineClientId { get; set; }
    
    /// <summary>
    /// Machine client associated with this token
    /// </summary>
    public virtual MachineClient? MachineClient { get; set; }
    
}