using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using follower_service.Models;
using follower_service.Services;
using System.Security.Claims;

namespace follower_service.Controllers;

[ApiController]
[Route("api/[controller]")]
public class FollowerController : ControllerBase
{
    private readonly Neo4jService _neo4jService;
    private readonly ILogger<FollowerController> _logger;

    public FollowerController(Neo4jService neo4jService, ILogger<FollowerController> logger)
    {
        _neo4jService = neo4jService;
        _logger = logger;
    }

    [HttpPost("follow")]
    [Authorize]
    public async Task<IActionResult> FollowUser([FromBody] FollowRequest request)
    {
        try
        {
            var userIdClaim = User.FindFirst(ClaimTypes.NameIdentifier)?.Value;
            if (string.IsNullOrEmpty(userIdClaim) || !int.TryParse(userIdClaim, out int followerId))
            {
                return Unauthorized(new { message = "Invalid user token" });
            }

            // Use authenticated user's ID as follower
            request.FollowerId = followerId;

            var success = await _neo4jService.FollowUserAsync(request.FollowerId, request.FollowingId);
            
            if (success)
            {
                return Ok(new { message = "User followed successfully" });
            }
            
            return BadRequest(new { message = "Cannot follow user" });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error in FollowUser endpoint");
            return StatusCode(500, new { message = "Internal server error" });
        }
    }

    [HttpDelete("unfollow/{followingId}")]
    [Authorize]
    public async Task<IActionResult> UnfollowUser(int followingId)
    {
        try
        {
            var userIdClaim = User.FindFirst(ClaimTypes.NameIdentifier)?.Value;
            _logger.LogInformation("UnfollowUser: userIdClaim = {Claim}, followingId = {FollowingId}", userIdClaim, followingId);
            
            if (string.IsNullOrEmpty(userIdClaim) || !int.TryParse(userIdClaim, out int followerId))
            {
                _logger.LogWarning("Invalid user token in UnfollowUser");
                return Unauthorized(new { message = "Invalid user token" });
            }

            _logger.LogInformation("Attempting to unfollow: follower={FollowerId}, following={FollowingId}", followerId, followingId);
            var success = await _neo4jService.UnfollowUserAsync(followerId, followingId);
            _logger.LogInformation("Unfollow result: {Success}", success);
            
            if (success)
            {
                return Ok(new { message = "User unfollowed successfully" });
            }
            
            _logger.LogWarning("Follow relationship not found between {FollowerId} and {FollowingId}", followerId, followingId);
            return NotFound(new { message = "Follow relationship not found" });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error in UnfollowUser endpoint");
            return StatusCode(500, new { message = "Internal server error" });
        }
    }

    [HttpGet("is-following/{followingId}")]
    [Authorize]
    public async Task<IActionResult> IsFollowing(int followingId)
    {
        try
        {
            var userIdClaim = User.FindFirst(ClaimTypes.NameIdentifier)?.Value;
            if (string.IsNullOrEmpty(userIdClaim) || !int.TryParse(userIdClaim, out int followerId))
            {
                return Unauthorized(new { message = "Invalid user token" });
            }

            var isFollowing = await _neo4jService.IsFollowingAsync(followerId, followingId);
            return Ok(new IsFollowingResponse { IsFollowing = isFollowing });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error in IsFollowing endpoint");
            return StatusCode(500, new { message = "Internal server error" });
        }
    }

    // Public endpoint - može se zvati bez autentifikacije (za SAGA pattern iz blog servisa)
    [HttpGet("check-following/{followerId}/{followingId}")]
    public async Task<IActionResult> CheckFollowing(int followerId, int followingId)
    {
        try
        {
            var isFollowing = await _neo4jService.IsFollowingAsync(followerId, followingId);
            return Ok(new IsFollowingResponse { IsFollowing = isFollowing });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error in CheckFollowing endpoint");
            return StatusCode(500, new { message = "Internal server error" });
        }
    }

    [HttpGet("followers")]
    [Authorize]
    public async Task<IActionResult> GetFollowers()
    {
        try
        {
            var userIdClaim = User.FindFirst(ClaimTypes.NameIdentifier)?.Value;
            if (string.IsNullOrEmpty(userIdClaim) || !int.TryParse(userIdClaim, out int userId))
            {
                return Unauthorized(new { message = "Invalid user token" });
            }

            var followers = await _neo4jService.GetFollowersAsync(userId);
            return Ok(followers);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error in GetFollowers endpoint");
            return StatusCode(500, new { message = "Internal server error" });
        }
    }

    [HttpGet("following")]
    [Authorize]
    public async Task<IActionResult> GetFollowing()
    {
        try
        {
            var userIdClaim = User.FindFirst(ClaimTypes.NameIdentifier)?.Value 
                           ?? User.FindFirst("sub")?.Value 
                           ?? User.FindFirst("userId")?.Value;
            
            _logger.LogInformation("GetFollowing: userIdClaim = {Claim}", userIdClaim);
            
            if (string.IsNullOrEmpty(userIdClaim) || !int.TryParse(userIdClaim, out int userId))
            {
                _logger.LogWarning("Invalid user token - userIdClaim: {Claim}", userIdClaim);
                return Unauthorized(new { message = "Invalid user token" });
            }

            var following = await _neo4jService.GetFollowingAsync(userId);
            _logger.LogInformation("User {UserId} is following {Count} users", userId, following.Count);
            return Ok(following);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error in GetFollowing endpoint");
            return StatusCode(500, new { message = "Internal server error" });
        }
    }

    [HttpGet("recommendations")]
    [Authorize]
    public async Task<IActionResult> GetRecommendations()
    {
        try
        {
            var userIdClaim = User.FindFirst(ClaimTypes.NameIdentifier)?.Value;
            if (string.IsNullOrEmpty(userIdClaim) || !int.TryParse(userIdClaim, out int userId))
            {
                return Unauthorized(new { message = "Invalid user token" });
            }

            var recommendations = await _neo4jService.GetRecommendationsAsync(userId);
            _logger.LogInformation("Found {Count} recommendations for user {UserId}", recommendations.Count, userId);
            return Ok(recommendations);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error in GetRecommendations endpoint");
            return StatusCode(500, new { message = "Internal server error" });
        }
    }

    [HttpGet("health")]
    public IActionResult HealthCheck()
    {
        return Ok(new { status = "healthy", service = "follower-service" });
    }
}
