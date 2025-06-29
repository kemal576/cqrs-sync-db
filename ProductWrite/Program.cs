using FluentValidation;
using Microsoft.EntityFrameworkCore;
using ProductWrite.Api;
using ProductWrite.Api.Handlers;
using ProductWrite.Application.Commands;
using ProductWrite.Domain.Base;
using ProductWrite.Infrastructure;
using Scalar.AspNetCore;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddOpenApi();

var mongoSettings = builder.Configuration.GetSection("MongoSettings").Get<MongoSettings>();
if (mongoSettings == null)
    throw new ArgumentNullException(nameof(mongoSettings));

// Service Registrations
builder.Services.AddDbContext<ApplicationDbContext>(options =>
    options.UseMongoDB(mongoSettings.ConnectionString, mongoSettings.DatabaseName));

builder.Services.AddMediatR(cfg => cfg.RegisterServicesFromAssembly(typeof(Program).Assembly));
builder.Services.AddValidatorsFromAssemblyContaining<CreateProductCommand>();
builder.Services.AddScoped<IProductRepository, ProductRepository>();
builder.Services.AddScoped<IUnitOfWork, ApplicationUnitOfWork>();


var app = builder.Build();

if (app.Environment.IsDevelopment())
{
    app.MapOpenApi();
    app.MapScalarApiReference();
}

app.MapProductEndpoints(); // handler registration
app.UseHttpsRedirection();

app.Run();