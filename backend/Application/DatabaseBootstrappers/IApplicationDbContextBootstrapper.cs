namespace Application.DatabaseBootstrappers;

public interface IApplicationDbContextBootstrapper
{
    Task BootstrapAsync(CancellationToken cancellationToken = default);
    
    void Bootstrap();
}