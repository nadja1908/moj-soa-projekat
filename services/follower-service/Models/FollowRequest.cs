namespace follower_service.Models;

public class FollowRequest
{
    public int FollowerId { get; set; }
    public int FollowingId { get; set; }
}

public class FollowerResponse
{
    public int UserId { get; set; }
    public string? Username { get; set; }
}

public class IsFollowingResponse
{
    public bool IsFollowing { get; set; }
}
