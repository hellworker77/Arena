using Domain.Entities.Abstract;

namespace Domain.Entities;

public class Tag : TimeLoggedEntity
{
    public string Name { get; set; }
    
    public string Description { get; set; }
}