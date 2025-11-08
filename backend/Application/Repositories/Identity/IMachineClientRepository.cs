using Domain.Entities.Identity;

namespace Application.Repositories.Identity;

public interface IMachineClientRepository
{
    Task<MachineClient?> GetMachineClientAsync(string clientId,
        CancellationToken cancellationToken);
}