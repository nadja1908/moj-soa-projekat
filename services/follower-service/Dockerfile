FROM mcr.microsoft.com/dotnet/aspnet:9.0 AS base
WORKDIR /app
EXPOSE 8080
RUN apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/*

FROM mcr.microsoft.com/dotnet/sdk:9.0 AS build
WORKDIR /src
COPY ["follower-service.csproj", "./"]
RUN dotnet restore "follower-service.csproj"
COPY . .
RUN dotnet build "follower-service.csproj" -c Release -o /app/build

FROM build AS publish
RUN dotnet publish "follower-service.csproj" -c Release -o /app/publish

FROM base AS final
WORKDIR /app
COPY --from=publish /app/publish .
ENTRYPOINT ["dotnet", "follower-service.dll"]
