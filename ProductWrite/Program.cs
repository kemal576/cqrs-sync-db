using FluentValidation;
using MediatR;
using Microsoft.Extensions.Options;
using MongoDB.Driver;
using ProductWrite.Api;
using ProductWrite.Api.Handlers;
using ProductWrite.Api.Validators;
using ProductWrite.Application.DistributedLock;
using ProductWrite.Infrastructure;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddOpenApi();
builder.Services.AddSwaggerGen();

var mongoSettings = builder.Configuration.GetSection("MongoSettings").Get<MongoSettings>();
if (mongoSettings == null)
    throw new ArgumentNullException(nameof(mongoSettings));

// Mongo Registration
builder.Services.Configure<MongoSettings>(
    builder.Configuration.GetSection("MongoSettings"));

builder.Services.AddSingleton<IMongoClient>(sp =>
{
    var settings = sp.GetRequiredService<IOptions<MongoSettings>>().Value;
    return new MongoClient(settings.ConnectionString);
});

builder.Services.AddScoped(sp =>
{
    var settings = sp.GetRequiredService<IOptions<MongoSettings>>().Value;
    var client = sp.GetRequiredService<IMongoClient>();
    return client.GetDatabase(settings.DatabaseName);
});

builder.Services.AddTransient(typeof(IPipelineBehavior<,>), typeof(RedLockBehavior<,>));
builder.Services.RegisterRedLock();

builder.Services.AddMediatR(cfg => cfg.RegisterServicesFromAssembly(typeof(Program).Assembly));
builder.Services.AddValidatorsFromAssemblyContaining<CreateProductRequestValidator>();
builder.Services.AddScoped<IProductRepository, ProductRepository>();

var app = builder.Build();

if (app.Environment.IsDevelopment())
{
    app.UseStaticFiles();
    app.MapOpenApi();
    app.UseSwagger();
    app.UseSwaggerUI(c =>{ c.InjectStylesheet("/SwaggerDark.css"); });
}

app.MapGet("/health", () => Results.Ok(new { status = "healthy" }));

app.MapProductEndpoints(); // handler registration
app.Run();