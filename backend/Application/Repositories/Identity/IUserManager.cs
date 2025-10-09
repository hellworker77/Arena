using Domain.Entities;

namespace Application.Repositories.Identity;

public interface IUserManager
{
    Task<ApplicationUser?> GetByLoginAsync(string login,
        CancellationToken cancellationToken = default);
    
    Task<bool> ValidatePasswordAsync(ApplicationUser user,
        string password,
        CancellationToken cancellationToken = default);
}