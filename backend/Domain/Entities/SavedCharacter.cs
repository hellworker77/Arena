using Domain.Entities.Abstract;

namespace Domain.Entities;

public partial class SavedCharacter: BaseEntity
{
    /// <summary>
    /// Version for preventing lose items
    /// </summary>
    public ulong Version { get; set; }
    
    /// <summary>
    /// Last touch date time
    /// </summary>
    public DateTime TouchDateTime { get; set; }
    
    /// <summary>
    /// Blob name in storage
    /// </summary>
    public string BlobName { get; set; } = string.Empty;
}