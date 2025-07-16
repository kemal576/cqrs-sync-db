using MongoDB.Driver;
using ProductWrite.Domain;

namespace ProductWrite.Infrastructure;

public interface IProductRepository
{
    Task AddAsync(Product product, CancellationToken cancellationToken);
}

public class ProductRepository : IProductRepository
{
    private readonly IMongoCollection<Product> _products;

    public ProductRepository(IMongoDatabase database)
    {
        _products = database.GetCollection<Product>("Products");
    }
    
    public async Task AddAsync(Product product, CancellationToken cancellationToken)
    {
        await _products.InsertOneAsync(product, cancellationToken);
    }
}
