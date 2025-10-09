namespace Application.Repositories.Identity;

public interface IUserRoleManager
{
    Task AddRoleToUserAsync(Guid userId,
        Guid roleId,
        CancellationToken cancellationToken = default);

    Task RemoveRoleFromUserAsync(Guid userId,
        Guid roleId,
        CancellationToken cancellationToken = default);
}