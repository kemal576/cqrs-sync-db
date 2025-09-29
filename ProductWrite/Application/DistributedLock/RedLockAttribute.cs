namespace ProductWrite.Application.DistributedLock;

[AttributeUsage(AttributeTargets.Class, Inherited = false)]
public sealed class RedLockAttribute(string resourceTemplate, int expirySeconds = 10) : Attribute
{
    private string ResourceTemplate { get; } = resourceTemplate;
    public int ExpirySeconds { get; } = expirySeconds;

    public string ResolveResource(object command)
    {
        var resource = ResourceTemplate;

        var props = command.GetType().GetProperties();
        foreach (var prop in props)
        {
            var value = prop.GetValue(command)?.ToString();
            resource = resource.Replace($"({prop.Name})", value ?? string.Empty);
        }

        return resource;
    }
}