using MongoDB.Driver;
using ProductWrite.Domain;

namespace ProductWrite.Infrastructure;

public interface IProductRepository
{
    Task<Product> GetByIdAsync(Guid id, CancellationToken cancellationToken);
    Task AddAsync(Product product, CancellationToken cancellationToken);
    Task UpdateAsync(Product product, CancellationToken cancellationToken);
}

public class ProductRepository : IProductRepository
{
    private readonly IMongoCollection<Product> _products;

    public ProductRepository(IMongoDatabase database)
    {
        _products = database.GetCollection<Product>("Products");
    }

    public async Task<Product> GetByIdAsync(Guid id, CancellationToken cancellationToken)
    {
        return await _products.Find(p => p.Id == id).FirstOrDefaultAsync(cancellationToken);
    }
    
    public async Task AddAsync(Product product, CancellationToken cancellationToken)
    {
        await _products.InsertOneAsync(product, new InsertOneOptions(), cancellationToken);
    }

    public async Task UpdateAsync(Product product, CancellationToken cancellationToken)
    {
        var update = Builders<Product>.Update.Set(p => p.Price, product.Price);
        await _products.UpdateOneAsync(p => p.Id == product.Id, update, cancellationToken: cancellationToken);
    }
}
