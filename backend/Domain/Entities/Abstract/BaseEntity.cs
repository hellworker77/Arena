namespace Domain.Entities.Abstract;

public abstract class BaseEntity
{
    /// <summary>
    /// Unique identifier of the entity
    /// </summary>
    public Guid Id { get; set; }
}