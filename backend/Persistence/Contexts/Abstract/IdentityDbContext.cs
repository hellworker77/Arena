using Domain.Entities.Identity;
using Domain.Entities.Identity.Abstract;
using Domain.Entities.Identity.Abstract.Relations;
using Microsoft.EntityFrameworkCore;

namespace Persistence.Contexts.Abstract;

public abstract class IdentityDbContext<TContext, TUser, TRole, TUserRole>(DbContextOptions<TContext> options)
    : BaseDbContext<TContext>(options)
    where TUserRole : BaseUserRole<TUser, TRole>
    where TUser : BaseUser<TUserRole>
    where TRole : BaseRole<TUserRole>
    where TContext : DbContext
{
    public DbSet<TUser> Users => Set<TUser>();

    public DbSet<TRole> Roles => Set<TRole>();

    public DbSet<TUserRole> UserRoles => Set<TUserRole>();

    public DbSet<JwtToken> Tokens => Set<JwtToken>();
    
    public DbSet<MachineClient> MachineClients => Set<MachineClient>();
    
    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        base.OnModelCreating(modelBuilder);

        modelBuilder.Entity<TUser>(entity =>
        {
            entity.HasKey(e => e.Id);
            entity.HasIndex(e => e.UserName);
            entity.HasIndex(e => e.Email).IsUnique();
            entity.Property(e => e.UserName).IsRequired().HasMaxLength(256);
            entity.Property(e => e.Email).IsRequired().HasMaxLength(256);
            entity.Property(e => e.PasswordHash).IsRequired();
        });

        modelBuilder.Entity<TRole>(entity =>
        {
            entity.HasKey(e => e.Id);
            entity.HasIndex(e => e.Name).IsUnique();
            entity.Property(e => e.Name).IsRequired().HasMaxLength(256);
        });

        modelBuilder.Entity<TUserRole>(entity =>
        {
            entity.HasKey(e => e.Id);

            entity.HasOne(ur => ur.User)
                .WithMany(u => u.UserRoles)
                .HasForeignKey(ur => ur.UserId)
                .OnDelete(DeleteBehavior.Cascade)
                .IsRequired();

            entity.HasOne(ur => ur.Role)
                .WithMany(r => r.UserRoles)
                .HasForeignKey(ur => ur.RoleId)
                .OnDelete(DeleteBehavior.Cascade)
                .IsRequired();
        });

        modelBuilder.Entity<JwtToken>(entity =>
        {
            entity.HasKey(t => t.Id);

            entity.HasOne(t => t.User)
                .WithMany(u => u.Tokens)
                .HasForeignKey(t => t.UserId)
                .OnDelete(DeleteBehavior.Cascade);

            entity.HasOne(t => t.MachineClient)
                .WithMany(m => m.Tokens)
                .HasForeignKey(t => t.MachineClientId)
                .OnDelete(DeleteBehavior.Cascade);
        });
    }
}