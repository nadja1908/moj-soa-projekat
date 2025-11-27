# 🔍 Observability Setup - Tracing, Logging & Monitoring

Ovaj dokument opisuje kompletno podešavanje observability stack-a za mikroservisnu aplikaciju koja uključuje **distributed tracing**, **centralized logging**, i **comprehensive monitoring**.

## 📊 Komponente Observability Stack-a

### 🎯 Distributed Tracing
- **Jaeger**: Prikupljanje i vizualizacija trace-ova
  - UI dostupan na: http://localhost:16686
  - Automatsko instrumentation za Go servise pomoću OpenTelemetry

### 📝 Centralized Logging  
- **Elasticsearch**: Storage i indeksiranje logova
- **Logstash**: Processing i parsing logova
- **Kibana**: Vizualizacija i pretražavanje logova
  - Kibana dostupna na: http://localhost:5601

### 📈 Monitoring & Metrics
- **Prometheus**: Metrics collection
  - Dostupan na: http://localhost:9090
- **Grafana**: Dashboard i vizualizacija metrika
  - **PAŽNJA**: Grafana je na portu **3001** umesto 3000 (da se izbegne konflikt)
  - Dostupna na: http://localhost:3001
  - Login: admin/admin
- **Node Exporter**: Host machine metrike
- **cAdvisor**: Container metrike

## 🚀 Pokretanje Observability Stack-a

```bash
# Pokreni celu aplikaciju sa observability
docker-compose up -d

# Proverava status svih servisa
docker-compose ps

# Za praćenje logova
docker-compose logs -f jaeger prometheus grafana elasticsearch kibana logstash
```

## 🔧 Portovi i Pristup

| Servis | Port | URL | Opis |
|--------|------|-----|------|
| **Jaeger UI** | 16686 | http://localhost:16686 | Distributed tracing |
| **Grafana** | 3001 | http://localhost:3001 | Metrics dashboards |
| **Prometheus** | 9090 | http://localhost:9090 | Metrics storage |
| **Kibana** | 5601 | http://localhost:5601 | Log analysis |
| **Elasticsearch** | 9200 | http://localhost:9200 | Log storage |
| **Node Exporter** | 9100 | http://localhost:9100/metrics | Host metrics |
| **cAdvisor** | 8081 | http://localhost:8081 | Container metrics |

## 📊 Dashboard-ovi u Grafani

### 1. Sistema i Kontejner Metrike
Prikazuje:
- CPU iskorišćenost host mašine
- Memoriju host mašine  
- Disk I/O operacije
- Mrežni saobraćaj
- CPU i memoriju po kontejneru

### 2. Mikroservisi - Performanse
Prikazuje:
- HTTP zahteve po minutu
- Response time (50th i 95th percentile)
- HTTP greške
- Go Garbage Collection metrike
- Goroutine broj

## 🔍 Implementirana Funkcionalnost

### ✅ Distributed Tracing
- **Auth Service**: Kompletno implementiran OpenTelemetry
- Automatski trace za sve HTTP zahteve
- Manuelni span-ovi za kritične operacije
- Trace context propagation između servisa

### ✅ Structured Logging
- **JSON format** logova za lakše parsiranje
- Integracija sa **Logstash** za centralno prikupljanje
- **Trace ID** i **Span ID** u svakom log entry
- Log nivoi: DEBUG, INFO, WARN, ERROR

### ✅ Metrics Collection
- **HTTP metrics**: Request count, duration, status codes
- **System metrics**: CPU, Memory, Disk, Network
- **Container metrics**: Resource usage po kontejneru  
- **Application metrics**: Custom business metrike

## 🛠️ Konfiguracija po Servisu

### Auth Service (Go)
```go
// Tracing
tracer := otel.Tracer(\"auth-service\")
_, span := tracer.Start(ctx, \"operation_name\")
defer span.End()

// Logging  
logrus.WithFields(logrus.Fields{
    \"username\": username,
    \"trace_id\": span.SpanContext().TraceID().String(),
}).Info(\"User login attempt\")

// Metrics endpoint
GET /metrics - Prometheus format
```

### Environment Variables
```bash
# Tracing
JAEGER_ENDPOINT=http://jaeger:14268/api/traces

# Logging
LOGSTASH_HOST=logstash:5000
```

## 📋 Zahtevi Projekta - Status

### ✅ Tracing Requirements
- [x] Implementiran distributed tracing u auth-service
- [x] Vizualizacija u Jaeger-u
- [x] OpenTelemetry instrumentation
- [x] Trace context propagation

### ✅ Logging Requirements  
- [x] Centralized log aggregation (ELK stack)
- [x] Structured JSON logging
- [x] Log processing u Logstash-u
- [x] Vizualizacija u Kibana-i

### ✅ Metrics Requirements
#### Host Machine Metrics:
- [x] CPU iskorišćenost
- [x] RAM memorija
- [x] File sistem
- [x] Mrežni saobraćaj

#### Container Metrics:
- [x] CPU iskorišćenost po kontejneru
- [x] RAM memorija po kontejneru  
- [x] File sistem po kontejneru
- [x] Mrežni saobraćaj po kontejneru

## 🔧 Dodavanje Observability u Druge Servise

Za dodavanje tracing-a u ostale Go servise:

1. **Dodaj dependencies u go.mod:**
```go
go.opentelemetry.io/otel v1.21.0
go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.46.1
github.com/sirupsen/logrus v1.9.3
```

2. **Kopiraj telemetry package** iz auth-service

3. **Inicijalizuj u main.go:**
```go
// Initialize tracing
tp, err := telemetry.InitTracer(serviceName, jaegerEndpoint)
defer telemetry.Shutdown(ctx, tp)

// Add middleware
r.Use(otelgin.Middleware(serviceName))
r.Use(telemetry.PrometheusMetrics())
```

4. **Dodaj environment variables** u docker-compose.yml

## 🐛 Troubleshooting

### Problem: Servis ne šalje trace-ove
**Rešenje**: Proverite da li je Jaeger accessible:
```bash
curl http://localhost:14268/api/traces
```

### Problem: Logovi se ne pojavljuju u Kibana
**Rešenje**: 
1. Proverite Elasticsearch: http://localhost:9200/_cat/indices
2. Proverite Logstash logs: `docker-compose logs logstash`

### Problem: Metrike se ne prikazuju u Grafana
**Rešenje**: Proverite Prometheus targets: http://localhost:9090/targets

## 🎯 Sledeći Koraci

1. **Dodaj tracing u ostale servise** (blog-service, tour-service, itd.)
2. **Konfiguriši alerting** u Grafana za kritične metrike  
3. **Dodaj business metrics** specifične za mikroservise
4. **Implementiraj log retention policy**
5. **Dodaj custom dashboard-ove** za različite scenarije

---

**🔥 NAPOMENA**: Grafana je na portu **3001** umesto 3000 da se izbegne konflikt sa postojećim servisima!