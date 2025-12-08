using Application.Repositories.Identity;
using Domain.Entities;
using Domain.Entities.Identity;
using Microsoft.EntityFrameworkCore;
using Persistence.Contexts;

namespace Persistence.Repositories.Identity;

public class MachineClientRepository(ApplicationDbContext dbContext): IMachineClientRepository
{
    public async Task<MachineClient?> GetMachineClientAsync(string clientId, CancellationToken cancellationToken)
        => await dbContext.MachineClients.FirstOrDefaultAsync(c => c.ClientId == clientId, cancellationToken);
}