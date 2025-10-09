using Domain.Entities.Abstract;

namespace Domain.Entities;

public class JwtToken: EventInfEntity
{
    public Guid Id { get; set; }
    
    /// <summary>
    /// Hash of the access token for security purposes
    /// </summary>
    public string AccessTokenHash { get; set; }
    
    public DateTime AccessTokenExpiresAt { get; set; }
    
    public bool IsRevoked { get; set; }
    
    /// <summary>
    /// Refresh token hash for security purposes
    /// </summary>
    public string RefreshTokenHash { get; set; }
    
    public DateTime RefreshTokenExpiresAt { get; set; }
    
    /// <summary>
    /// Foreign key to the user
    /// </summary>
    public Guid UserId { get; set; }
    
    public virtual ApplicationUser User { get; set; }
}