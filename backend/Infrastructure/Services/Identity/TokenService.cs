using System.IdentityModel.Tokens.Jwt;
using System.Security.Claims;
using System.Security.Cryptography;
using System.Text;
using Application.Contracts.Data.Identity;
using Application.Repositories.Identity;
using Application.Services.Identity;
using Domain.Entities;
using Microsoft.Extensions.Configuration;
using Microsoft.IdentityModel.Tokens;
using Persistence.Contexts;

namespace Infrastructure.Services.Identity;

public class TokenService : ITokenService
{
    private readonly RSA _privateRsa;
    private readonly IConfiguration _configuration;
    private readonly ApplicationDbContext _dbContext;
    private readonly RsaPem _pem;
    private readonly IUserManager _userManager;
    private readonly IRoleManager _roleManager;
    private readonly IJwtTokenRepository _jwtTokenRepository;

    private static readonly int AccessTokenLifetimeMinutes = 60;
    private static readonly int RefreshTokenLifetimeDays = 7;

    public TokenService(ApplicationDbContext dbContext,
        IConfiguration configuration,
        RsaPem pem, 
        IUserManager userManager,
        IRoleManager roleManager,
        IJwtTokenRepository jwtTokenRepository)
    {
        _pem = pem;
        _userManager = userManager;
        _roleManager = roleManager;
        _jwtTokenRepository = jwtTokenRepository;
        _dbContext = dbContext;
        _configuration = configuration;

        _privateRsa = RSA.Create();
        _privateRsa.ImportFromPem(_pem.privatePem);
    }

    public async Task<Jwks> GetJwksAsync(CancellationToken cancellationToken = default)
    {
        using var rsaPublic = RSA.Create();
        
        rsaPublic.ImportFromPem(_pem.publicPem);

        var key = new RsaSecurityKey(rsaPublic);
        var jwk = JsonWebKeyConverter.ConvertFromRSASecurityKey(key);
        jwk.Kid = "main-key";

        var result = new Jwks
        {
            Keys =
            [
                new()
                {
                    Kid = jwk.Kid,
                    N = jwk.N!,
                    E = jwk.E!
                }
            ]
        };
        
        return await Task.FromResult(result);
    }

    public async Task<Jwt> CreateJwtAsync(BaseLoginDto loginDto,
        CancellationToken cancellationToken = default)
    {
        var user = await _userManager.GetByLoginAsync(loginDto.Login, cancellationToken)
            ?? throw new UnauthorizedAccessException("Invalid login or password");

        if (!await _userManager.ValidatePasswordAsync(user, loginDto.Password, cancellationToken))
            throw new UnauthorizedAccessException("Invalid login or password");

        var accessToken = await GenerateJwtAsync(user, cancellationToken);
        var refreshToken = GenerateRefreshToken();
        var hashedRefreshToken = HashToken(refreshToken);

        var accessLifeTime = _configuration.GetValue("Jwt:AccessTokenLifetimeMinutes", AccessTokenLifetimeMinutes);
        var refreshLifeTime = _configuration.GetValue("Jwt:RefreshTokenLifetimeDays", RefreshTokenLifetimeDays);
        
        var jwt = new JwtToken
        {
            AccessTokenHash = accessToken,
            AccessTokenExpiresAt = DateTime.UtcNow.AddMinutes(accessLifeTime),
            RefreshTokenHash = hashedRefreshToken,
            RefreshTokenExpiresAt = DateTime.UtcNow.AddDays(refreshLifeTime),
            UserId = user.Id,
            IsRevoked = false
        };

        await _jwtTokenRepository.SaveJwtAsync(jwt, cancellationToken);
        await _jwtTokenRepository.RetainOldTokensAsync(user.Id, cancellationToken: cancellationToken);
        
        return new Jwt
        {
            AccessToken = accessToken,
            RefreshToken = refreshToken,
            ExpiresAt = jwt.AccessTokenExpiresAt
        };
    }

    public async Task<Jwt> RefreshJwtAsync(string refreshToken,
        CancellationToken cancellationToken = default)
    {
        var hashedToken = HashToken(refreshToken);

        var tokenEntity = await _jwtTokenRepository.GetByRefreshTokenHashAsync(hashedToken, cancellationToken)
                          ?? throw new UnauthorizedAccessException("Invalid or expired refresh token");
        
        if (tokenEntity.IsRevoked || tokenEntity.RefreshTokenExpiresAt <= DateTime.UtcNow)
            throw new UnauthorizedAccessException("Invalid or expired refresh token");

        tokenEntity.IsRevoked = true;
        await _jwtTokenRepository.UpdateJwtAsync(tokenEntity, cancellationToken);
        
        var newRefreshToken = GenerateRefreshToken();
        var newHashedRefreshToken = HashToken(newRefreshToken);
        
        var accessLifeTime = _configuration.GetValue("Jwt:AccessTokenLifetimeMinutes", AccessTokenLifetimeMinutes);
        var refreshLifeTime = _configuration.GetValue("Jwt:RefreshTokenLifetimeDays", RefreshTokenLifetimeDays);
        
        var newToken  = new JwtToken
        {
            AccessTokenHash = tokenEntity.AccessTokenHash,
            AccessTokenExpiresAt = DateTime.UtcNow.AddMinutes(accessLifeTime),
            RefreshTokenHash = newHashedRefreshToken,
            RefreshTokenExpiresAt = DateTime.UtcNow.AddDays(refreshLifeTime),
            UserId = tokenEntity.UserId,
            IsRevoked = false
        };
        
        await _jwtTokenRepository.SaveJwtAsync(newToken, cancellationToken);
        await _jwtTokenRepository.RetainOldTokensAsync(tokenEntity.UserId, cancellationToken: cancellationToken);
        
        await _dbContext.SaveChangesAsync(cancellationToken);

        return new Jwt
        {
            AccessToken = newToken.AccessTokenHash,
            RefreshToken = newRefreshToken,
            ExpiresAt =  newToken.AccessTokenExpiresAt
        };
    }

    public async Task RevokeRefreshTokenAsync(string refreshToken,
        CancellationToken cancellationToken = default)
    {
        var hashed = HashToken(refreshToken);
        
        var tokenEntity = await _jwtTokenRepository.GetByRefreshTokenHashAsync(hashed, cancellationToken);
        
        if (tokenEntity != null)
        {
            tokenEntity.IsRevoked = true;
            await _jwtTokenRepository.UpdateJwtAsync(tokenEntity, cancellationToken);
        }
    }

    private async Task<string> GenerateJwtAsync(ApplicationUser user,
        CancellationToken cancellationToken = default)
    {
        var credentials = new SigningCredentials(new RsaSecurityKey(_privateRsa), SecurityAlgorithms.RsaSha256);

        var roleNames = await _roleManager.GetUserRolesAsync(user.Id, cancellationToken);

        var claims = new List<Claim>
        {
            new(JwtRegisteredClaimNames.Sub, user.Id.ToString()),
            new(JwtRegisteredClaimNames.Name, user.UserName),
            new(JwtRegisteredClaimNames.Email, user.Email),
            new(JwtRegisteredClaimNames.Jti, Guid.NewGuid().ToString())
        };
        
        claims.AddRange(roleNames.Select(rn => new Claim(ClaimTypes.Role, rn)));

        var accessLifeTime = _configuration.GetValue("Jwt:AccessTokenLifetimeMinutes", AccessTokenLifetimeMinutes);
        
        var token = new JwtSecurityToken(
            issuer: _configuration["Jwt:Issuer"] ?? "AuthService",
            audience: "Microservices",
            claims: claims,
            expires: DateTime.UtcNow.AddMinutes(accessLifeTime),
            signingCredentials: credentials);

        return new JwtSecurityTokenHandler().WriteToken(token);
    }

    private static string GenerateRefreshToken() =>
        Convert.ToBase64String(RandomNumberGenerator.GetBytes(64));
    
    private static string HashToken(string token)
    {
        using var sha = SHA256.Create();
        return Convert.ToHexString(sha.ComputeHash(Encoding.UTF8.GetBytes(token)));
    }
}