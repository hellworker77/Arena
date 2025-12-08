using Application.Repositories.Identity;
using Domain.Entities.Identity;
using Microsoft.EntityFrameworkCore;
using Persistence.Contexts;

namespace Persistence.Repositories.Identity;

public class UserManager(ApplicationDbContext dbContext): IUserManager
{
    public async Task<ApplicationUser?> GetByLoginAsync(string login,
        CancellationToken cancellationToken = default)
        => await dbContext.Users
            .FirstOrDefaultAsync(u => u.UserName == login || u.Email == login, cancellationToken);

    public Task<bool> ValidatePasswordAsync(ApplicationUser user,
        string password,
        CancellationToken cancellationToken = default)
    {
        var valid = BCrypt.Net.BCrypt.Verify(password, user.PasswordHash);
        return Task.FromResult(valid);
    }
}