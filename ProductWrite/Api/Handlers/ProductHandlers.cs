using FluentValidation;
using MediatR;
using ProductWrite.Application.Commands;

namespace ProductWrite.Api.Handlers;

public static class ProductHandlers
{
    public static void MapProductEndpoints(this IEndpointRouteBuilder routes)
    {
        routes.MapPost("/products", async (
            CreateProductCommand command,
            IMediator mediator,
            IValidator<CreateProductCommand> validator) =>
        {
            var validation = await validator.ValidateAsync(command);
            if (!validation.IsValid)
                return Results.BadRequest(validation.Errors.Select(e => e.ErrorMessage));

            var productId = await mediator.Send(command);
            
            return Results.Ok(new { productId });
        });
    }
}