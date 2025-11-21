using Neo4j.Driver;

namespace follower_service.Services;

public class Neo4jService : IDisposable
{
    private readonly IDriver _driver;
    private readonly ILogger<Neo4jService> _logger;

    public Neo4jService(IConfiguration configuration, ILogger<Neo4jService> logger)
    {
        _logger = logger;
        var uri = configuration["Neo4j:Uri"] ?? "bolt://neo4j:7687";
        var user = configuration["Neo4j:User"] ?? "neo4j";
        var password = configuration["Neo4j:Password"] ?? "password";

        _driver = GraphDatabase.Driver(uri, AuthTokens.Basic(user, password));
        _logger.LogInformation("Neo4j driver created successfully");
    }

    public async Task<bool> FollowUserAsync(int followerId, int followingId)
    {
        if (followerId == followingId)
        {
            _logger.LogWarning("User {FollowerId} attempted to follow themselves", followerId);
            return false;
        }

        await using var session = _driver.AsyncSession();
        try
        {
            var result = await session.ExecuteWriteAsync(async tx =>
            {
                var query = @"
                    MERGE (follower:User {id: $followerId})
                    MERGE (following:User {id: $followingId})
                    MERGE (follower)-[r:FOLLOWS]->(following)
                    RETURN r";

                var cursor = await tx.RunAsync(query, new { followerId, followingId });
                var record = await cursor.SingleAsync();
                return record != null;
            });

            _logger.LogInformation("User {FollowerId} followed user {FollowingId}", followerId, followingId);
            return result;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error following user {FollowingId} by user {FollowerId}", followingId, followerId);
            throw;
        }
    }

    public async Task<bool> UnfollowUserAsync(int followerId, int followingId)
    {
        await using var session = _driver.AsyncSession();
        try
        {
            var result = await session.ExecuteWriteAsync(async tx =>
            {
                var query = @"
                    MATCH (follower:User {id: $followerId})-[r:FOLLOWS]->(following:User {id: $followingId})
                    DELETE r
                    RETURN count(r) as deleted";

                var cursor = await tx.RunAsync(query, new { followerId, followingId });
                var record = await cursor.SingleAsync();
                return record["deleted"].As<int>() > 0;
            });

            _logger.LogInformation("User {FollowerId} unfollowed user {FollowingId}", followerId, followingId);
            return result;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error unfollowing user {FollowingId} by user {FollowerId}", followingId, followerId);
            throw;
        }
    }

    public async Task<bool> IsFollowingAsync(int followerId, int followingId)
    {
        await using var session = _driver.AsyncSession();
        try
        {
            var result = await session.ExecuteReadAsync(async tx =>
            {
                var query = @"
                    MATCH (follower:User {id: $followerId})-[:FOLLOWS]->(following:User {id: $followingId})
                    RETURN count(*) > 0 as isFollowing";

                var cursor = await tx.RunAsync(query, new { followerId, followingId });
                var record = await cursor.SingleAsync();
                return record["isFollowing"].As<bool>();
            });

            return result;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error checking if user {FollowerId} follows user {FollowingId}", followerId, followingId);
            return false;
        }
    }

    public async Task<List<int>> GetFollowersAsync(int userId)
    {
        await using var session = _driver.AsyncSession();
        try
        {
            var followers = await session.ExecuteReadAsync(async tx =>
            {
                var query = @"
                    MATCH (follower:User)-[:FOLLOWS]->(user:User {id: $userId})
                    RETURN follower.id as followerId";

                var cursor = await tx.RunAsync(query, new { userId });
                var records = await cursor.ToListAsync();
                return records.Select(r => r["followerId"].As<int>()).ToList();
            });

            return followers;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error getting followers for user {UserId}", userId);
            throw;
        }
    }

    public async Task<List<int>> GetFollowingAsync(int userId)
    {
        await using var session = _driver.AsyncSession();
        try
        {
            var following = await session.ExecuteReadAsync(async tx =>
            {
                var query = @"
                    MATCH (user:User {id: $userId})-[:FOLLOWS]->(following:User)
                    RETURN following.id as followingId";

                var cursor = await tx.RunAsync(query, new { userId });
                var records = await cursor.ToListAsync();
                return records.Select(r => r["followingId"].As<int>()).ToList();
            });

            return following;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error getting following for user {UserId}", userId);
            throw;
        }
    }

    public void Dispose()
    {
        _driver?.Dispose();
    }
}
