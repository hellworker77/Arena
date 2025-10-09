namespace Application.Contracts.Data.Identity;

public class Jwks
{
    public record JwkKey
    {
        public string Kty { get; set; } = "RSA";
        public string Kid { get; set; } = null!;
        public string Use { get; set; } = "sig";
        public string Alg { get; set; } = "RS256";
        public string N { get; set; } = null!;
        public string E { get; set; } = null!;
    }
    
    public List<JwkKey> Keys { get; set; } = new();
}