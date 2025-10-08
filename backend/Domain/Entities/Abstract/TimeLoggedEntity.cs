namespace Domain.Entities.Abstract;

public abstract class TimeLoggedEntity : BaseEntity
{
    public DateTime CreatedAt { get; set; }
    
    public DateTime LastUpdateAt { get; set; }
}