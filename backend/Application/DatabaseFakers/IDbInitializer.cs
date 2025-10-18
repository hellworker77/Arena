namespace Application.DatabaseFakers;

public interface IDbInitializer
{
    Task InitializeAsync(CancellationToken cancellationToken = default);
}