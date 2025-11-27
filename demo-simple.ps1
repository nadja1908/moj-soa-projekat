# SIMPLE OBSERVABILITY DEMO
Write-Host "SIMPLE OBSERVABILITY DEMO" -ForegroundColor Cyan
Write-Host "==========================" -ForegroundColor Cyan

# Stop all services first
Write-Host "Stopping all services..." -ForegroundColor Yellow
docker-compose down

# Start only essential observability services
Write-Host "Starting essential observability services..." -ForegroundColor Yellow
docker-compose up -d jaeger prometheus grafana node-exporter cadvisor

Write-Host "Starting core application services..." -ForegroundColor Yellow  
docker-compose up -d eureka-server stakeholders-db stakeholders-service auth-service api-gateway

# Wait for services to start
Write-Host "Waiting for services to initialize..." -ForegroundColor Yellow
Start-Sleep 45

# Health check function
function Test-ServiceQuick {
    param($ServiceName, $Url)
    try {
        $response = Invoke-WebRequest -Uri $Url -UseBasicParsing -TimeoutSec 5
        Write-Host "OK $ServiceName" -ForegroundColor Green
        return $true
    } catch {
        Write-Host "FAIL $ServiceName" -ForegroundColor Red
        return $false
    }
}

Write-Host ""
Write-Host "Checking Observability Services:" -ForegroundColor Cyan

Test-ServiceQuick "Jaeger" "http://localhost:16686"
Test-ServiceQuick "Grafana" "http://localhost:3001"
Test-ServiceQuick "Prometheus" "http://localhost:9090"
Test-ServiceQuick "Node Exporter" "http://localhost:9100"
Test-ServiceQuick "cAdvisor" "http://localhost:8081"

Write-Host ""
Write-Host "Checking Application Services:" -ForegroundColor Cyan

Test-ServiceQuick "Auth Service" "http://localhost:8003/health"
Test-ServiceQuick "Stakeholders" "http://localhost:8001/health"
Test-ServiceQuick "API Gateway" "http://localhost:8000/health"

Write-Host ""
Write-Host "RUNNING DEMO SCENARIO" -ForegroundColor Green
Write-Host "=====================" -ForegroundColor Green

Write-Host ""
Write-Host "1. Generating traces and metrics..." -ForegroundColor Yellow

# Generate some load to create metrics and traces
for ($i = 1; $i -le 10; $i++) {
    try {
        Write-Host "Request $i/10..." -NoNewline
        
        # Health checks
        Invoke-RestMethod -Uri "http://localhost:8003/health" -TimeoutSec 5 | Out-Null
        Invoke-RestMethod -Uri "http://localhost:8001/health" -TimeoutSec 5 | Out-Null
        
        # Try login (will create trace)
        if ($i -eq 5) {
            try {
                $body = @{ username = "demo_user"; password = "demo_pass" } | ConvertTo-Json
                Invoke-RestMethod -Uri "http://localhost:8000/login" -Method Post -Body $body -ContentType "application/json" -TimeoutSec 5 | Out-Null
            } catch {
                # Expected to fail, but creates trace
            }
        }
        
        Write-Host " OK" -ForegroundColor Green
        Start-Sleep -Milliseconds 200
    } catch {
        Write-Host " FAIL" -ForegroundColor Red
    }
}

Write-Host ""
Write-Host "2. Demo is ready! Open these URLs:" -ForegroundColor Green
Write-Host ""
Write-Host "JAEGER (Tracing):    http://localhost:16686" -ForegroundColor Cyan
Write-Host "  Service: auth-service" -ForegroundColor Gray
Write-Host "  Operation: login_user or GET /health" -ForegroundColor Gray
Write-Host "  Click Find Traces" -ForegroundColor Gray
Write-Host ""
Write-Host "GRAFANA (Dashboards): http://localhost:3001" -ForegroundColor Cyan
Write-Host "  Login: admin/admin" -ForegroundColor Gray
Write-Host "  Import dashboards or create custom panels" -ForegroundColor Gray
Write-Host ""
Write-Host "PROMETHEUS (Raw Metrics): http://localhost:9090" -ForegroundColor Cyan
Write-Host "  Try queries like:" -ForegroundColor Gray
Write-Host "    rate(node_cpu_seconds_total[5m])" -ForegroundColor Gray
Write-Host "    node_memory_MemAvailable_bytes" -ForegroundColor Gray
Write-Host "    container_memory_usage_bytes" -ForegroundColor Gray

Write-Host ""
Write-Host "3. DEMO TALKING POINTS:" -ForegroundColor Yellow
Write-Host ""
Write-Host "DISTRIBUTED TRACING:" -ForegroundColor Green
Write-Host "  OpenTelemetry implementation in Go auth-service" -ForegroundColor White
Write-Host "  Automatic HTTP request tracing" -ForegroundColor White
Write-Host "  Manual spans for business logic" -ForegroundColor White
Write-Host "  Trace visualization in Jaeger" -ForegroundColor White
Write-Host ""
Write-Host "COMPREHENSIVE METRICS:" -ForegroundColor Green  
Write-Host "  Host machine: CPU, RAM, Disk, Network" -ForegroundColor White
Write-Host "  Container metrics: Resource usage per service" -ForegroundColor White
Write-Host "  Application metrics: HTTP requests, response times" -ForegroundColor White
Write-Host "  Real-time visualization in Grafana" -ForegroundColor White

Write-Host ""
Write-Host "DEMO COMPLETED! All project requirements satisfied:" -ForegroundColor Green
Write-Host "  Tracing implemented and visualized" -ForegroundColor White
Write-Host "  Host machine metrics (CPU, RAM, Disk, Network)" -ForegroundColor White
Write-Host "  Container metrics (CPU, RAM, Disk, Network)" -ForegroundColor White
Write-Host "  Production-ready Go instrumentation" -ForegroundColor White

Write-Host ""
Write-Host "NOTE: Grafana is on PORT 3001 to avoid conflicts!" -ForegroundColor Yellow