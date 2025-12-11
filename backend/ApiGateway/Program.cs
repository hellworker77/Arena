using Application.DatabaseBootstrappers;
using Application.DatabaseFakers;
using Microsoft.AspNetCore.Authentication.JwtBearer;
using Ocelot.DependencyInjection;
using Ocelot.Middleware;

var builder = WebApplication.CreateBuilder(args);

builder.Configuration.AddJsonFile("ocelot.json", optional: false, reloadOnChange: true);
builder.Services.AddOcelot(builder.Configuration);

builder.Services
    .AddAuthentication(JwtBearerDefaults.AuthenticationScheme)
    .AddJwtBearer("JwtBearer", options =>
    {
        options.Authority = builder.Configuration["Identity:Authority"];
        options.Audience = builder.Configuration["Identity:Audience"];
        options.RequireHttpsMetadata = false;
    });

var app = builder.Build();

using (var scope = app.Services.CreateScope())
{
    /*var bootstrapper = scope.ServiceProvider.GetRequiredService<IApplicationDbContextBootstrapper>();
    await bootstrapper.BootstrapAsync();
    
    var initializer = scope.ServiceProvider.GetRequiredService<IDbInitializer>();
    await initializer.InitializeAsync();*/
}

app.UseAuthentication();
app.UseAuthorization();

await app.UseOcelot();

app.Run();