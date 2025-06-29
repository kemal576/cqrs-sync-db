using MongoDB.Bson.Serialization.Attributes;

namespace ProductWrite.Domain.Base;

public class BaseEntity
{
    [BsonId]
    public Guid Id { get; init; }
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
}