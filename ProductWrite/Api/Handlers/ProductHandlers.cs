using FluentValidation;
using MediatR;
using ProductWrite.Application.Commands;

namespace ProductWrite.Api.Handlers;

public static class ProductHandlers
{
    public static void MapProductEndpoints(this IEndpointRouteBuilder routes)
    {
        routes.MapPost("/products", async (
            Requests.CreateProductRequest request,
            IMediator mediator,
            IValidator<Requests.CreateProductRequest> validator) =>
        {
            var validation = await validator.ValidateAsync(request);
            if (!validation.IsValid)
                return Results.BadRequest(validation.Errors.Select(e => e.ErrorMessage));

            var command = new CreateProductCommand(new Application.Inputs.CreateProductCommandInput
            {
                Name = request.Name,
                Description = request.Description,
                Price = request.Price
            });

            var productId = await mediator.Send(command);

            return Results.Ok(new { productId });
        });

        routes.MapPatch("/products/{id}/price", async (
            Guid id,
            Requests.UpdatePriceRequest request,
            IMediator mediator,
            IValidator<UpdateProductPriceCommand> validator) =>
        {
            var command = new UpdateProductPriceCommand(id, request.Price);

            var validation = await validator.ValidateAsync(command);
            if (!validation.IsValid)
                return Results.BadRequest(validation.Errors.Select(e => e.ErrorMessage));

            await mediator.Send(command);

            return Results.NoContent();
        });
    }
}