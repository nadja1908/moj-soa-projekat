# Observability Stack - Health Check Script (PowerShell)
Write-Host "🔍 Observability Stack - Health Check Script" -ForegroundColor Cyan
Write-Host "===============================================" -ForegroundColor Cyan

function Test-Service {
    param(
        [string]$ServiceName,
        [string]$Url,
        [int]$ExpectedCode = 200
    )
    
    Write-Host "Checking $ServiceName... " -NoNewline
    
    try {
        $response = Invoke-WebRequest -Uri $Url -Method Get -UseBasicParsing -TimeoutSec 10
        if ($response.StatusCode -eq $ExpectedCode) {
            Write-Host "✓ OK" -ForegroundColor Green
            return $true
        } else {
            Write-Host "✗ FAILED (Status: $($response.StatusCode))" -ForegroundColor Red
            return $false
        }
    } catch {
        # For services like Grafana that redirect (302), check if it's a redirect
        if ($_.Exception.Response.StatusCode.value__ -eq $ExpectedCode) {
            Write-Host "✓ OK (Redirect)" -ForegroundColor Green
            return $true
        } else {
            Write-Host "✗ FAILED ($($_.Exception.Message))" -ForegroundColor Red
            return $false
        }
    }
}

Write-Host ""
Write-Host "🎯 Observability Services:" -ForegroundColor Yellow
Write-Host "-------------------------"

Test-Service "Jaeger UI" "http://localhost:16686"
Test-Service "Grafana" "http://localhost:3001" 302
Test-Service "Prometheus" "http://localhost:9090"
Test-Service "Kibana" "http://localhost:5601" 302
Test-Service "Elasticsearch" "http://localhost:9200"
Test-Service "Node Exporter" "http://localhost:9100"
Test-Service "cAdvisor" "http://localhost:8081"

Write-Host ""
Write-Host "📊 Application Services:" -ForegroundColor Yellow
Write-Host "------------------------"

Test-Service "Auth Service" "http://localhost:8003/health"
Test-Service "Auth Metrics" "http://localhost:8003/metrics"
Test-Service "API Gateway" "http://localhost:8000/health"
Test-Service "Stakeholders Service" "http://localhost:8001/health"

Write-Host ""
Write-Host "🧪 Testing Tracing:" -ForegroundColor Yellow
Write-Host "-------------------"

Write-Host "Sending test login request to generate traces..."
try {
    $body = @{
        username = "testuser"
        password = "testpass"
    } | ConvertTo-Json

    $response = Invoke-RestMethod -Uri "http://localhost:8000/login" -Method Post -Body $body -ContentType "application/json" -TimeoutSec 10
    Write-Host "✓ Test request sent" -ForegroundColor Green
    Write-Host "ℹ Check traces in Jaeger: http://localhost:16686" -ForegroundColor Yellow
} catch {
    Write-Host "✗ Failed to send test request: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""
Write-Host "📈 Quick Access URLs:" -ForegroundColor Yellow
Write-Host "---------------------"
Write-Host "🎯 Jaeger (Tracing):     http://localhost:16686"
Write-Host "📊 Grafana (Dashboards): http://localhost:3001 (admin/admin)"
Write-Host "⚡ Prometheus (Metrics):  http://localhost:9090"
Write-Host "🔍 Kibana (Logs):        http://localhost:5601"
Write-Host ""
Write-Host "Note: Grafana is on port 3001 to avoid conflicts!" -ForegroundColor Yellow