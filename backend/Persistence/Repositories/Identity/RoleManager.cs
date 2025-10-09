using Application.Repositories.Identity;
using Microsoft.EntityFrameworkCore;
using Persistence.Contexts;

namespace Persistence.Repositories.Identity;

public class RoleManager(ApplicationDbContext dbContext): IRoleManager
{
    public async Task<IReadOnlyList<string>> GetUserRolesAsync(Guid userId,
        CancellationToken cancellationToken = default)
        => await dbContext.Roles
            .Include(r => r.UserRoles)
            .Where(r => r.UserRoles.Any(ur => ur.UserId == userId))
            .Select(r => r.Name)
            .ToListAsync(cancellationToken);
}