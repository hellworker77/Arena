using Combat;
using Infrastructure.Hubs;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddCors(options =>
{
    options.AddPolicy("AllowFrontend", policy =>
    {
        policy
            .WithOrigins("http://localhost:5173") 
            .AllowAnyHeader()
            .AllowAnyMethod()
            .AllowCredentials(); 
    });
});

builder.Services
    .AddServices()
    .AddRepositories()
    .AddFactories()
    .AddWebRtc()
    .AddFabrics()
    .AddEventBus()
    .AddStores();
    
builder.Services.AddControllers();
builder.Services.AddSignalR();

var app = builder.Build();

app.UseCors("AllowFrontend");

app.MapControllers();
app.MapHub<ArenaSignalHub>("/arenaSignalHub");
app.Run();
