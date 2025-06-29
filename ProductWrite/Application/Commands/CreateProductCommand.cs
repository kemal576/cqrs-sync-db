using MediatR;
using MongoDB.Bson;
using ProductWrite.Application.Inputs;
using ProductWrite.Domain;
using ProductWrite.Domain.Base;
using ProductWrite.Infrastructure;

namespace ProductWrite.Application.Commands;

public record CreateProductCommand(CreateProductCommandInput Input) : IRequest<Guid>;

public class CreateProductCommandHandler : IRequestHandler<CreateProductCommand, Guid>
{
    private readonly IProductRepository _repository;
    private readonly IUnitOfWork _unitOfWork;

    public CreateProductCommandHandler(IProductRepository repository, IUnitOfWork unitOfWork)
    {
        _repository = repository;
        _unitOfWork = unitOfWork;
    }

    public async Task<Guid> Handle(CreateProductCommand request,
        CancellationToken cancellationToken)
    {
        var product = Product.Create(request.Input.Name, request.Input.Description, request.Input.Price);

        await _repository.AddAsync(product,  cancellationToken);
        await _unitOfWork.SaveChangesAsync(cancellationToken);

        return product.Id;
    }
}