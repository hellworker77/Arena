using Domain.Entities.Abstract;

namespace Domain.Entities;

public class Media: TimeLoggedEntity
{
    public string Name { get; set; }
    
    public Guid BlobId { get; set; }
    
    public string ContentType { get; set; }
    
    public long Size { get; set; }
}