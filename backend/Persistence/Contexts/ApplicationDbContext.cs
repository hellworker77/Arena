using Domain.Entities;
using Domain.Entities.Relations;
using Persistence.Contexts.Abstract;

namespace Persistence.Contexts;

public class ApplicationDbContext 
    : IdentityDbContext<ApplicationUser, ApplicationRole, ApplicationUserRole>
{
    
}