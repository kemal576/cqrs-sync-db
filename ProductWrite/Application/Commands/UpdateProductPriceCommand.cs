using FluentValidation;
using MediatR;
using ProductWrite.Infrastructure;

namespace ProductWrite.Application.Commands;

public record UpdateProductPriceCommand(Guid Id, decimal Price) : IRequest;

public class UpdateProductPriceCommandHandler : IRequestHandler<UpdateProductPriceCommand>
{
    private readonly IProductRepository _repository;

    public UpdateProductPriceCommandHandler(IProductRepository repository)
    {
        _repository = repository;
    }

    public async Task Handle(UpdateProductPriceCommand request, CancellationToken cancellationToken)
    {
        var product = await _repository.GetByIdAsync(request.Id, cancellationToken);
        if (product == null)
            throw new KeyNotFoundException($"Product with ID {request.Id} not found.");

        product.UpdatePrice(request.Price);
        await _repository.UpdateAsync(product, cancellationToken);
    }
}
