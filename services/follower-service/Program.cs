using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.IdentityModel.Tokens;
using System.Text;
using follower_service.Services;
using Prometheus;
using System.Diagnostics;
using OpenTelemetry.Resources;
using OpenTelemetry.Trace;
using OpenTelemetry.Exporter;
using Serilog;
using Serilog.Formatting.Compact;

// Configure Serilog for JSON logging to Console ONLY (Docker će prikupiti logove)
Log.Logger = new LoggerConfiguration()
    .WriteTo.Console(new CompactJsonFormatter())
    .MinimumLevel.Information()
    .CreateLogger();

var builder = WebApplication.CreateBuilder(args);

// Add Serilog
builder.Host.UseSerilog();

// Add OpenTelemetry Tracing
builder.Services.AddOpenTelemetry()
    .WithTracing(tracing => tracing
        .SetResourceBuilder(ResourceBuilder.CreateDefault()
            .AddService("follower-service", "1.0.0"))
        .AddAspNetCoreInstrumentation()
        .AddHttpClientInstrumentation()
        .AddConsoleExporter()
        .AddOtlpExporter(options =>
        {
            options.Endpoint = new Uri("http://jaeger:4317");
            options.Protocol = OtlpExportProtocol.Grpc;
        }));

// Add services to the container
builder.Services.AddControllers();
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddOpenApi();

// Add HTTP request logging
builder.Services.AddHttpLogging(logging =>
{
    logging.LoggingFields = Microsoft.AspNetCore.HttpLogging.HttpLoggingFields.RequestPropertiesAndHeaders |
                           Microsoft.AspNetCore.HttpLogging.HttpLoggingFields.ResponsePropertiesAndHeaders;
    logging.RequestHeaders.Clear();
    logging.ResponseHeaders.Clear();
    logging.MediaTypeOptions.Clear();
});

// Add Neo4j Service
builder.Services.AddSingleton<Neo4jService>();

// Add JWT Authentication
var jwtSecret = builder.Configuration["Jwt:Secret"] ?? "your-secret-key-min-16-chars-long-for-security";
builder.Services.AddAuthentication(JwtBearerDefaults.AuthenticationScheme)
    .AddJwtBearer(options =>
    {
        options.TokenValidationParameters = new TokenValidationParameters
        {
            ValidateIssuerSigningKey = true,
            IssuerSigningKey = new SymmetricSecurityKey(Encoding.UTF8.GetBytes(jwtSecret)),
            ValidateIssuer = false,
            ValidateAudience = false,
            ValidateLifetime = true,
            ClockSkew = TimeSpan.Zero
        };
    });

builder.Services.AddAuthorization();

// Add CORS
builder.Services.AddCors(options =>
{
    options.AddPolicy("AllowAll", policy =>
    {
        policy.AllowAnyOrigin()
              .AllowAnyMethod()
              .AllowAnyHeader();
    });
});

var app = builder.Build();

// Configure the HTTP request pipeline
if (app.Environment.IsDevelopment())
{
    app.MapOpenApi();
}

app.UseCors("AllowAll");

// Add HTTP request logging
app.UseHttpLogging();

// Prometheus metrics middleware
app.UseRouting();
app.UseHttpMetrics();
app.MapMetrics();

app.UseAuthentication();
app.UseAuthorization();
app.MapControllers();

app.Run();
