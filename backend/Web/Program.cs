using FelCache;
using Web;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddFelCache(options =>
{
    options.Ttl = 30000;
    
    options.BaseUrl = builder.Configuration["FelCache:Url"];
});

builder.Services.AddFelAuthentication(builder.Configuration);

builder.Services.AddOpenApi();

builder.Services.AddControllers();

builder.Services.AddEndpointsApiExplorer();

var app = builder.Build();

app.UseAuthentication();
app.UseAuthorization();

app.UseHttpsRedirection();

app.MapControllers();

app.Run();