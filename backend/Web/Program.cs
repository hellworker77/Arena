using Application.DatabaseBootstrappers;
using Application.DatabaseFakers;
using FelCache;
using Web;

var builder = WebApplication.CreateBuilder(args);

builder.Services
    .AddFelCache(options =>
    {
        options.Ttl = 30000;

        options.BaseUrl = builder.Configuration["FelCache:Url"];
    })
    .AddFelAuthentication(builder.Configuration)
    .AddDbContexts(builder.Configuration);

builder.Services.AddControllers();

var app = builder.Build();

app.UseAuthentication();
app.UseAuthorization();

app.MapControllers();

app.Run();