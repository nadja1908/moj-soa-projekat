namespace follower_service.Models;

public class UserRecommendation
{
    public int UserId { get; set; }
    public string Username { get; set; } = string.Empty;
    public string Email { get; set; } = string.Empty;
    public string? FirstName { get; set; }
    public string? LastName { get; set; }
    public string? Role { get; set; }
    public int CommonFollowers { get; set; }  // Broj zajedničkih pratioca
}
