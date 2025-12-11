namespace Constants;

public static class CharacterConstants
{
    
    public static Func<string, ulong, string> GenerateBlobName = (userId, version) =>
        $"{userId}/saved_characters/character_v{version}.sav";
    
    public const byte MaxSavedCharactersPerUser = 5;
    
    public const byte MaxVersionsPerCharacter = 50;
}