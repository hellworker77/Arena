using Application.Repositories.Identity;
using Domain.Entities;
using Microsoft.EntityFrameworkCore;
using Persistence.Contexts;

namespace Persistence.Repositories.Identity;

public class JwtTokenRepository(ApplicationDbContext dbContext) : IJwtTokenRepository
{
    public async Task<JwtToken?> GetByRefreshTokenHashAsync(string refreshTokenHash,
        CancellationToken cancellationToken = default)
        => await dbContext.Tokens
            .Include(t => t.User)
            .FirstOrDefaultAsync(t => t.RefreshTokenHash == refreshTokenHash, cancellationToken);
    
    public async Task SaveJwtAsync(JwtToken token,
        CancellationToken cancellationToken = default)
    {
        dbContext.Tokens.Add(token);
        await dbContext.SaveChangesAsync(cancellationToken);
    }

    public async Task UpdateJwtAsync(JwtToken token,
        CancellationToken cancellationToken = default)
    {
        dbContext.Tokens.Update(token);
        await dbContext.SaveChangesAsync(cancellationToken);
    }

    public async Task<int> RetainOldTokensAsync(Guid userId,
        int keepCount = 5,
        CancellationToken cancellationToken = default)
    {
        var tokens = await dbContext.Tokens
            .Where(t => t.UserId == userId)
            .OrderByDescending(t => t.CreatedDate)
            .ToListAsync(cancellationToken);
        
        return await RetainOldTokensAsync(tokens, keepCount, cancellationToken);
    }

    public async Task<int> RetainOldMachineClientTokensAsync(Guid machineClientId, int keepCount = 5,
        CancellationToken cancellationToken = default)
    {
        var tokens = await dbContext.Tokens
            .Where(t => t.MachineClientId == machineClientId)
            .OrderByDescending(t => t.CreatedDate)
            .ToListAsync(cancellationToken);

        return await RetainOldTokensAsync(tokens, keepCount, cancellationToken);
    }

    private async Task<int> RetainOldTokensAsync(List<JwtToken> tokens,
        int keepCount,
        CancellationToken cancellationToken = default)
    {
        if (tokens.Count <= keepCount)
            return 0;

        var oldTokens = tokens
            .Skip(keepCount)
            .ToList();
        
        dbContext.Tokens.RemoveRange(oldTokens);
        await dbContext.SaveChangesAsync(cancellationToken);
        return oldTokens.Count;
    }
}