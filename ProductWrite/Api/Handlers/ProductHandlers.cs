using FluentValidation;
using MediatR;
using ProductWrite.Api.Requests;
using ProductWrite.Application.Commands;

namespace ProductWrite.Api.Handlers;

public static class ProductHandlers
{
    public static void MapProductEndpoints(this IEndpointRouteBuilder routes)
    {
        routes.MapPost("/products", async (
            CreateProductRequest request,
            IMediator mediator,
            IValidator<CreateProductRequest> validator) =>
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

        routes.MapPatch("/products/{id:guid}/price", async (
            Guid id,
            UpdateProductPriceRequest request,
            IMediator mediator,
            IValidator<UpdateProductPriceRequest> validator) =>
        {
            var validation = await validator.ValidateAsync(request);
            if (!validation.IsValid)
                return Results.BadRequest(validation.Errors.Select(e => e.ErrorMessage));
            
            var command = new UpdateProductPriceCommand(id, request.Price);

            await mediator.Send(command);

            return Results.NoContent();
        });
    }
}