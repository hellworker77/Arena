using Identity;

var builder = WebApplication.CreateBuilder(args);

builder.Services
    .AddRsaPem()
    .AddRepositories()
    .AddServices()
    .AddDbContexts(builder.Configuration);

builder.Services.AddControllers();

var app = builder.Build();

app.MapControllers();

app.Run();