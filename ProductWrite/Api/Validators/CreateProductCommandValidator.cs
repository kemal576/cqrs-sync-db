using FluentValidation;
using ProductWrite.Application.Commands;

namespace ProductWrite.Api.Validators;

public class CreateProductCommandValidator : AbstractValidator<CreateProductCommand>
{
    public CreateProductCommandValidator()
    {
        RuleFor(x => x.Input.Name).NotEmpty().WithMessage("Name cannot be null or empty");
        RuleFor(x => x.Input.Description).NotEmpty().WithMessage("Description cannot be null or empty");
        RuleFor(x => x.Input.Price).GreaterThan(0).WithMessage("Price cannot zero or negative");
    }
}