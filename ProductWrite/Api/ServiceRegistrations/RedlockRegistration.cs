using RedLockNet.SERedis;
using RedLockNet.SERedis.Configuration;
using StackExchange.Redis;

public static class RedLockRegistration
{
    public static IServiceCollection RegisterRedLock(this IServiceCollection services)
    {
        var redisNodesEnv = Environment.GetEnvironmentVariable("REDIS_NODES") ?? "";
        var redisNodes = redisNodesEnv.Split(',', StringSplitOptions.RemoveEmptyEntries | StringSplitOptions.TrimEntries);

        var multiplexers = redisNodes
            .Select(node => ConnectionMultiplexer.Connect(node))
            .Select(conn => new RedLockMultiplexer(conn))
            .ToList();

        var redlockFactory = RedLockFactory.Create(multiplexers);
        services.AddSingleton(redlockFactory);

        return services;
    }
}
