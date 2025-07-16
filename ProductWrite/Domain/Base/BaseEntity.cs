using MongoDB.Bson;
using MongoDB.Bson.Serialization.Attributes;

namespace ProductWrite.Domain.Base;

public class BaseEntity
{
    [BsonId]
    [BsonRepresentation(BsonType.String)]
    public Guid Id { get; protected init; } = Guid.NewGuid();

    public DateTime CreatedAt { get; protected set; } = DateTime.UtcNow;
    public DateTime UpdatedAt { get; protected set; } = DateTime.UtcNow;

    public void TouchUpdatedAt() => UpdatedAt = DateTime.UtcNow;
}