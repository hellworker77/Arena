namespace Domain.Entities.Abstract;

public abstract class EventInfEntity : BaseEntity
{
    /// <summary>
    /// Creation date of the entity
    /// </summary>
    public DateTime CreatedDate { get; set; }
    
    /// <summary>
    /// Last modification date of the entity
    /// </summary>
    public DateTime LastModifiedDate { get; set; }
}