using Application.DatabaseFakers;
using Domain.Entities.Identity;
using Domain.Entities.Identity.Relations;
using Persistence.Contexts;

namespace Persistence.DatabaseFakers;

public class DbInitializer(ApplicationDbContext dbContext) : IDbInitializer
{
    public async Task InitializeAsync(CancellationToken cancellationToken = default)
    {
        await dbContext.Database.EnsureDeletedAsync(cancellationToken);
        await dbContext.Database.EnsureCreatedAsync(cancellationToken);
        
        dbContext.Users.AddRange(FakeData.Users);
        await dbContext.SaveChangesAsync(cancellationToken);
        
        dbContext.Roles.AddRange(FakeData.Roles);
        await dbContext.SaveChangesAsync(cancellationToken);
        
        dbContext.UserRoles.AddRange(FakeData.userRoles);
        await dbContext.SaveChangesAsync(cancellationToken);
        
        dbContext.Tokens.AddRange(FakeData.jwtTokens);
        await dbContext.SaveChangesAsync(cancellationToken);
    }

    private record FakeData
    {
        
        public static readonly IEnumerable<ApplicationUser> Users = new List<ApplicationUser>()
        {
            new()
            {
                Id = Guid.Parse("11111111-1111-1111-1111-111111111111"),
                UserName = "admin",
                Email = "admin@gmain.com",
                PasswordHash = BCrypt.Net.BCrypt.HashPassword("Admin@123"),
            }
        };
        
        public static readonly IEnumerable<ApplicationRole> Roles = new List<ApplicationRole>()
        {
            new ()
            {
                Id = Guid.Parse("22222222-2222-2222-2222-222222222222"),
                Name = "Admin"
            }
        };
        
        public static readonly IEnumerable<ApplicationUserRole> userRoles = new List<ApplicationUserRole>()
        {
            new ()
            {
                UserId = Guid.Parse("11111111-1111-1111-1111-111111111111"),
                RoleId = Guid.Parse("22222222-2222-2222-2222-222222222222")
            }
        };
        
        public static readonly IEnumerable<JwtToken> jwtTokens = new List<JwtToken>();
    }
}