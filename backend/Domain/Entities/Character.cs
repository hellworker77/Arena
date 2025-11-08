using Domain.Contracts;
using Domain.Entities.Abstract;

namespace Domain.Entities;

public partial class Character: EventInfEntity
{
    public string Name { get; set; } = string.Empty;
    
    public CharacterClass Class { get; set; }
}