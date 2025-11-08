namespace Application.Contracts.Data.Identity;

public class Jwt
{
    public string? AccessToken { get; set; }

    public string RefreshToken { get; set; } = string.Empty;
    
    public DateTime ExpiresAt { get; set; }
}