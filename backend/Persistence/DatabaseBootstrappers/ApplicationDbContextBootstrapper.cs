using Application.DatabaseBootstrappers;
using Microsoft.EntityFrameworkCore;
using Persistence.Contexts;

namespace Persistence.DatabaseBootstrappers;

public class ApplicationDbContextBootstrapper(ApplicationDbContext dbContext) : IApplicationDbContextBootstrapper
{
    public async Task BootstrapAsync(CancellationToken cancellationToken = default)
    {
        var pending = await dbContext.Database.GetPendingMigrationsAsync(cancellationToken);
        
        if (!pending.Any())
            return;
        
        await dbContext.Database.MigrateAsync(cancellationToken);
    }

    public void Bootstrap()
    {
        var pending = dbContext.Database.GetPendingMigrations();
        
        if (!pending.Any())
            return;
        
        dbContext.Database.Migrate();
    }
}