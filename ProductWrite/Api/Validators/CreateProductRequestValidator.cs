using FluentValidation;
using ProductWrite.Api.Requests;

namespace ProductWrite.Api.Validators;

public class CreateProductRequestValidator : AbstractValidator<CreateProductRequest>
{
    public CreateProductRequestValidator()
    {
        RuleFor(x => x.Name).NotEmpty().WithMessage("Name cannot be null or empty");
        RuleFor(x => x.Description).NotEmpty().WithMessage("Description cannot be null or empty");
        RuleFor(x => x.Price).GreaterThan(0).WithMessage("Price cannot zero or negative");
    }
}
