namespace Domain.Entities.Abstract;

public class ApplicationUser: TimeLoggedEntity 
{
    public string UserName { get; set; }
    
    public string Email { get; set; }
    
    public string PasswordHash { get; set; }
}