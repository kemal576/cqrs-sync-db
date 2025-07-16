using MediatR;
using ProductWrite.Application.Inputs;
using ProductWrite.Domain;
using ProductWrite.Infrastructure;

namespace ProductWrite.Application.Commands;

public record CreateProductCommand(CreateProductCommandInput Input) : IRequest<Guid>;

public class CreateProductCommandHandler : IRequestHandler<CreateProductCommand, Guid>
{
    private readonly IProductRepository _repository;

    public CreateProductCommandHandler(IProductRepository repository)
    {
        _repository = repository;
    }

    public async Task<Guid> Handle(CreateProductCommand request,
        CancellationToken cancellationToken)
    {
        var product = Product.Create(request.Input.Name, request.Input.Description, request.Input.Price);

        await _repository.AddAsync(product,  cancellationToken);

        return product.Id;
    }
}