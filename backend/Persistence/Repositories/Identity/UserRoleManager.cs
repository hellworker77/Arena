using Application.Repositories.Identity;
using Domain.Entities.Identity.Relations;
using Microsoft.EntityFrameworkCore;
using Persistence.Contexts;

namespace Persistence.Repositories.Identity;

public class UserRoleManager(ApplicationDbContext dbContext) : IUserRoleManager
{
    public async Task AddRoleToUserAsync(Guid userId,
        Guid roleId,
        CancellationToken cancellationToken = default)
    {
        var relation = new ApplicationUserRole
        {
            UserId = userId,
            RoleId = roleId
        };

        dbContext.UserRoles.Add(relation);
        await dbContext.SaveChangesAsync(cancellationToken);
    }

    public async Task RemoveRoleFromUserAsync(Guid userId,
        Guid roleId,
        CancellationToken cancellationToken = default)
    {
        var relation = await dbContext.UserRoles
            .FirstOrDefaultAsync(ur => ur.UserId == userId && ur.RoleId == roleId, cancellationToken);
        if (relation != null)
        {
            dbContext.UserRoles.Remove(relation);
            await dbContext.SaveChangesAsync(cancellationToken);
        }
    }
}