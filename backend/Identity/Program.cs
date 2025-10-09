using Identity;

var builder = WebApplication.CreateBuilder(args);

builder.Services
    .AddRsaPem()
    .AddRepositories()
    .AddServices();

builder.Services.AddOpenApi();

builder.Services.AddControllers();

builder.Services.AddEndpointsApiExplorer();

var app = builder.Build();

app.UseHttpsRedirection();

app.MapControllers();

app.Run();