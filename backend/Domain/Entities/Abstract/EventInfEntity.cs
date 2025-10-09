namespace Domain.Entities.Abstract;

public abstract class EventInfEntity : BaseEntity
{
    public DateTime CreatedDate { get; set; }
    
    public DateTime LastModifiedDate { get; set; }
}