using QAService;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddControllers();

builder.Services
    .ConfigureCors(builder.Configuration)
    .AddFelAuthentication(builder.Configuration);

var app = builder.Build();

app.MapControllers();

app.Run();