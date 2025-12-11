using ApiGateway;
using Application.DatabaseBootstrappers;
using Application.DatabaseFakers;
using Ocelot.Middleware;

var builder = WebApplication.CreateBuilder(args);

builder.Services
    .AddFelAuthentication(builder.Configuration)
    .AddApiGateway(builder.Configuration, builder.Environment.ContentRootPath)
    .AddDbContexts(builder.Configuration);

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

await app.UseOcelot();

app.Run();