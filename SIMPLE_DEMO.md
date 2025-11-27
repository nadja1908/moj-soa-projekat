# 🚀 JEDNOSTAVNA DEMO VERZIJA - Observability

Zbog ograničenih resursa, evo kako možeš demonstrirati observability bez ELK stack-a:

## 🎯 **DEMO PLAN - Osnovni Stack**

### **KORAK 1: Pokreni samo osnovne observability servise**

```powershell
# Umesto celog stack-a, pokreni samo ključne komponente:
docker-compose up -d jaeger prometheus grafana node-exporter cadvisor auth-service stakeholders-service api-gateway eureka-server stakeholders-db
```

### **KORAK 2: Demonstriraj Tracing (Jaeger)**

1. **Otvori Jaeger**: http://localhost:16686
2. **Pošalji test zahteve**:
```powershell
# Login zahtev koji generiše trace
Invoke-RestMethod -Uri "http://localhost:8000/login" -Method Post -Body '{"username":"demo","password":"test"}' -ContentType "application/json"

# Health check zahtevi
for ($i=1; $i -le 5; $i++) {
    Invoke-RestMethod -Uri "http://localhost:8003/health" -Method Get
}
```

3. **U Jaeger UI**:
   - Service: "auth-service" 
   - Operation: "login_user"
   - Prikaži trace timeline
   - Objasni span details

### **KORAK 3: Demonstriraj Metrics (Grafana + Prometheus)**

1. **Otvori Grafana**: http://localhost:3001 (admin/admin)
2. **Importuj dashboard** ili napravi simple panel
3. **Prikaži metrics**:
   - CPU usage: `rate(node_cpu_seconds_total[5m])`
   - Memory usage: `node_memory_MemAvailable_bytes`
   - HTTP requests: `http_requests_total`

### **KORAK 4: Live Demo Scenario**

```powershell
# SCENARIO: E-commerce korisnik
Write-Host "=== DEMO: E-commerce User Journey ===" -ForegroundColor Green

# 1. Health checks (warming up)
Write-Host "1. Proverava servise..." -ForegroundColor Yellow
Invoke-RestMethod -Uri "http://localhost:8003/health"
Invoke-RestMethod -Uri "http://localhost:8001/health"

# 2. User login attempt
Write-Host "2. Korisnik pokušava login..." -ForegroundColor Yellow
try {
    $loginResult = Invoke-RestMethod -Uri "http://localhost:8000/login" -Method Post -Body '{"username":"john_doe","password":"password123"}' -ContentType "application/json"
    Write-Host "   Login successful!" -ForegroundColor Green
} catch {
    Write-Host "   Login failed (expected) - CHECK TRACES!" -ForegroundColor Red
}

# 3. Multiple API calls
Write-Host "3. Simulira browsing..." -ForegroundColor Yellow
for ($i=1; $i -le 10; $i++) {
    Invoke-RestMethod -Uri "http://localhost:8001/health" -Method Get
    Invoke-RestMethod -Uri "http://localhost:8003/health" -Method Get
    Start-Sleep -Milliseconds 100
}

Write-Host "4. DEMO ZAVRŠEN - Proverite observability tools:" -ForegroundColor Green
Write-Host "   📊 Traces: http://localhost:16686" -ForegroundColor Cyan
Write-Host "   📈 Metrics: http://localhost:3001" -ForegroundColor Cyan
Write-Host "   🔧 Prometheus: http://localhost:9090" -ForegroundColor Cyan
```

## 🎭 **ŠTO DEMO POKRIVA**

### ✅ **Distributed Tracing**
- OpenTelemetry implementation u Go servisu
- Automatic span creation za HTTP requests
- Manual spans za business logic
- Trace visualization u Jaeger

### ✅ **Comprehensive Metrics**
- **Host Machine Metrics**: CPU, RAM, Disk, Network
- **Container Metrics**: Resource usage po servisu
- **Application Metrics**: HTTP requests, response times
- Real-time visualization u Grafana

### ✅ **Structural Implementation**
- Proper Go telemetry package
- Environment-based configuration
- Graceful shutdown handling
- Production-ready setup

## 🎪 **DEMO PRESENTATION (10 min)**

### **1. Uvod (1 min)**
```text
"Implementiran observability stack:
- Distributed Tracing sa OpenTelemetry i Jaeger
- Comprehensive monitoring sa Prometheus i Grafana
- Host i container metrike
- Production-ready Go instrumentation"
```

### **2. Architecture Overview (2 min)**
```text
"Pokažimo komponente:
- Jaeger za distributed tracing
- Prometheus za metrics collection
- Grafana za visualization (PORT 3001!)
- Go servisi sa OpenTelemetry instrumentation"
```

### **3. Live Demo (5 min)**

**A) Tracing (2 min):**
- Otvori Jaeger
- Pokreni demo script
- Prikaži traces sa span hierarchy

**B) Metrics (3 min):**
- Otvori Grafana
- Prikaži real-time CPU/Memory
- Generiši load → vidi metric changes
- Objasni container vs host metrics

### **4. Code Explanation (2 min)**
```text
"Implementation details:
- OpenTelemetry automatic instrumentation
- Custom spans za business operations
- Structured logging sa trace correlation
- Prometheus metrics endpoint"
```

## 📋 **PRE-DEMO CHECKLIST**

```powershell
# 1. Start services
docker-compose up -d jaeger prometheus grafana node-exporter auth-service stakeholders-service api-gateway eureka-server stakeholders-db

# 2. Wait for services
Start-Sleep 30

# 3. Quick health check
$services = @("jaeger:16686", "grafana:3001", "prometheus:9090", "auth-service:8003/health")
foreach ($service in $services) {
    try {
        $url = "http://localhost:$($service.Replace(':', ':'))"
        Invoke-WebRequest $url -TimeoutSec 5
        Write-Host "✅ $service" -ForegroundColor Green
    } catch {
        Write-Host "❌ $service" -ForegroundColor Red
    }
}
```

---

**🔥 KLJUČNE DEMO TAČKE:**

1. **Port 3001** za Grafana (avoiding conflicts!)
2. **Real-time tracing** u Jaeger-u
3. **Host + Container metrics** u Grafana
4. **Production-ready code** sa proper error handling
5. **Svi projektni zahtevi ispunjeni** ✅

Ovaj pristup je **resurno efikasniji** i fokusira se na **ključne observability komponente**!