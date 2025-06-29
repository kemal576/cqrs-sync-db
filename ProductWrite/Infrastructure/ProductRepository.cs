using ProductWrite.Domain;

namespace ProductWrite.Infrastructure;

public interface IProductRepository
{
    Task AddAsync(Product product, CancellationToken cancellationToken);
}

public class ProductRepository(ApplicationDbContext context) : IProductRepository
{
    public async Task AddAsync(Product product, CancellationToken cancellationToken)
    {
        await context.Products.AddAsync(product, cancellationToken);
    }
}
