using ProductWrite.Domain.Base;

namespace ProductWrite.Domain;

public class Product : BaseEntity
{
    public string Name { get; private set; }
    public string Description { get; private set; }
    public decimal Price { get; private set; }

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
            Name = name,
            Description = description,
            Price = price,
        };
    }

    public void UpdatePrice(decimal price)
    {
        if (price <= 0)
            throw new ArgumentException("Price cannot be zero or negative.", nameof(price));

        Price = price;
    }
}

