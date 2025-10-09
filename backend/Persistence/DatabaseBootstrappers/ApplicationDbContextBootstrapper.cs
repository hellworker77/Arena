using Application.DatabaseBootstrappers;
using Microsoft.EntityFrameworkCore;
using Persistence.Contexts;

namespace Persistence.DatabaseBootstrappers;

public class ApplicationDbContextBootstrapper(ApplicationDbContext dbContext) : IApplicationDbContextBootstrapper
{
    public async Task BootstrapAsync(CancellationToken cancellationToken = default)
        => await dbContext.Database.MigrateAsync(cancellationToken);

    public void Bootstrap()
        => dbContext.Database.Migrate();
}