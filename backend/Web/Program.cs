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

using (var scope = app.Services.CreateScope())
{
    var bootstrapper = scope.ServiceProvider.GetRequiredService<IApplicationDbContextBootstrapper>();
    await bootstrapper.BootstrapAsync();
    
    var initializer = scope.ServiceProvider.GetRequiredService<IDbInitializer>();
    await initializer.InitializeAsync();
}

app.UseAuthentication();
app.UseAuthorization();

app.MapControllers();

app.Run();