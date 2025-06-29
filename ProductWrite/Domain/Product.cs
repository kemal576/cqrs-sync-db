using MongoDB.Bson;
using MongoDB.EntityFrameworkCore;
using ProductWrite.Domain.Base;

namespace ProductWrite.Domain;

[Collection("products")]
public class Product : BaseEntity
{
    public string Name {get; init;}
    public string Description {get; init;}
    public decimal Price {get; init;}

    private Product()
    {
    }

    public static Product Create(string name, string description, decimal price)
    {
        if (string.IsNullOrWhiteSpace(name))
            throw new ArgumentException("Name cannot be null or whitespace.", nameof(name));
        
        if (string.IsNullOrWhiteSpace(description))
            throw new ArgumentException("Description cannot be null or whitespace.", nameof(name));
        
        if (price <= 0)
            throw new ArgumentException("Price cannot be zero or negative.", nameof(price));
            
        return new Product
        { 
            Id = Guid.NewGuid(),
            Name = name,
            Description = description,
            Price = price,
        };
    }
}

