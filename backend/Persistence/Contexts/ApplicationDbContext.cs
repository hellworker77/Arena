using Domain.Entities.Identity;
using Domain.Entities.Identity.Relations;
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