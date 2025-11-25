using Neo4j.Driver;
using follower_service.Models;

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
                    MERGE (follower:User {userId: $followerId})
                    MERGE (following:User {userId: $followingId})
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
            var deletedCount = await session.ExecuteWriteAsync(async tx =>
            {
                var query = @"
                    MATCH (follower:User {userId: $followerId})-[r:FOLLOWS]->(following:User {userId: $followingId})
                    DELETE r
                    RETURN count(*) as deletedCount";

                var cursor = await tx.RunAsync(query, new { followerId, followingId });
                var summary = await cursor.ConsumeAsync();
                return summary.Counters.RelationshipsDeleted;
            });

            var success = deletedCount > 0;
            _logger.LogInformation("User {FollowerId} unfollowed user {FollowingId}, deleted {Count} relationships, success: {Result}", 
                followerId, followingId, deletedCount, success);
            return success;
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
                    MATCH (follower:User {userId: $followerId})-[:FOLLOWS]->(following:User {userId: $followingId})
                    RETURN COUNT(*) > 0 as isFollowing";

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
                    MATCH (follower:User)-[:FOLLOWS]->(user:User {userId: $userId})
                    RETURN follower.userId as followerId";

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
                    MATCH (user:User {userId: $userId})-[:FOLLOWS]->(following:User)
                    RETURN following.userId as followingId";

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

    public async Task<List<UserRecommendation>> GetRecommendationsAsync(int userId)
    {
        await using var session = _driver.AsyncSession();
        try
        {
            var recommendations = await session.ExecuteReadAsync(async tx =>
            {
                // Prvo pokušaj da nađeš prijatelje od prijatelja
                var friendsOfFriendsQuery = @"
                    MATCH (me:User {userId: $userId})-[:FOLLOWS]->(friend:User)-[:FOLLOWS]->(recommendation:User)
                    WHERE NOT (me)-[:FOLLOWS]->(recommendation) 
                      AND recommendation.userId <> $userId
                    WITH recommendation, count(DISTINCT friend) as commonFollowers
                    RETURN recommendation.userId as userId, 
                           recommendation.username as username,
                           commonFollowers
                    ORDER BY commonFollowers DESC
                    LIMIT 10";

                var cursor = await tx.RunAsync(friendsOfFriendsQuery, new { userId });
                var records = await cursor.ToListAsync();
                
                // Ako nema prijatelja od prijatelja, vrati sve korisnike koje ne pratiš
                if (records.Count == 0)
                {
                    var allUsersQuery = @"
                        MATCH (me:User {userId: $userId})
                        MATCH (recommendation:User)
                        WHERE recommendation.userId <> $userId
                          AND NOT (me)-[:FOLLOWS]->(recommendation)
                        RETURN recommendation.userId as userId,
                               recommendation.username as username,
                               0 as commonFollowers
                        LIMIT 10";
                    
                    cursor = await tx.RunAsync(allUsersQuery, new { userId });
                    records = await cursor.ToListAsync();
                }
                
                return records.Select(r => new UserRecommendation
                {
                    UserId = r["userId"].As<int>(),
                    Username = r["username"].As<string>(),
                    CommonFollowers = r["commonFollowers"].As<int>(),
                    Email = "",  // Email će se popuniti iz stakeholders servisa
                    Role = null
                }).ToList();
            });

            _logger.LogInformation("Found {Count} raw recommendations from Neo4j for user {UserId}", recommendations.Count, userId);
            
            // Sada treba da pozovemo stakeholders servis da dobijemo dodatne informacije
            // Za sada vraćamo samo osnovne informacije iz Neo4j
            return recommendations;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error getting recommendations for user {UserId}", userId);
            throw;
        }
    }

    public void Dispose()
    {
        _driver?.Dispose();
    }
}
