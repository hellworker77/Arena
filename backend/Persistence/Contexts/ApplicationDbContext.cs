using Domain.Entities;
using Domain.Entities.Relations;
using Microsoft.EntityFrameworkCore;
using Persistence.Contexts.Abstract;

namespace Persistence.Contexts;

public class ApplicationDbContext (DbContextOptions<ApplicationDbContext> options)
    : IdentityDbContext<ApplicationDbContext,
        ApplicationUser,
        ApplicationRole,
        ApplicationUserRole>(options)
{
    
}