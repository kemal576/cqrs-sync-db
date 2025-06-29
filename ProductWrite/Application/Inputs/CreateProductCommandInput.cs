namespace ProductWrite.Application.Inputs;

public class CreateProductCommandInput
{
    public string Name { get; init; }
    public string Description { get; init; }
    public decimal Price { get; init; }
}