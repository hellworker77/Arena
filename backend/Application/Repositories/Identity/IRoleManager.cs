namespace Application.Repositories.Identity;

public interface IRoleManager
{
    Task<IReadOnlyList<string>> GetUserRolesAsync(Guid userId,
        CancellationToken cancellationToken = default);
}