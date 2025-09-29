using MediatR;
using RedLockNet.SERedis;

namespace ProductWrite.Application.DistributedLock;

public class RedLockBehavior<TRequest, TResponse>(RedLockFactory redLockFactory)
    : IPipelineBehavior<TRequest, TResponse> where TRequest : notnull
{
    public async Task<TResponse> Handle(
        TRequest request,
        RequestHandlerDelegate<TResponse> next,
        CancellationToken cancellationToken)
    {
        if (typeof(TRequest).GetCustomAttributes(typeof(RedLockAttribute), false)
                .FirstOrDefault() is not RedLockAttribute attr)
        {
            return await next();
        }

        var resource = attr.ResolveResource(request);
        var expiry = TimeSpan.FromSeconds(attr.ExpirySeconds);

        // No retry on purpose. Code can be refactored easily with retry timespan parameter
        await using var redLock = await redLockFactory.CreateLockAsync(resource, expiry);
        if (!redLock.IsAcquired)
            throw new InvalidOperationException($"Could not acquire lock for resource '{resource}'");

        return await next();
    }
}

