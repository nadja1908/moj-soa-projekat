# RPC Protocol Implementation

Ova implementacija dodaje RPC (Remote Procedure Call) protokol između API Gateway-a i servisa, kao dodatna alternativa postojećoj HTTP REST komunikaciji.

## Implementirane RPC funkcionalnosti

### 1. Auth Service RPC
**Endpoint:** `auth-service:9003` (RPC port)

#### RPC metode:
- `AuthService.Login` - Autentifikacija korisnika
- `AuthService.Register` - Registracija novog korisnika

#### API Gateway rute:
- `POST /api/rpc/auth/login` - RPC login
- `POST /api/rpc/auth/register` - RPC registracija

### 2. Tour Service RPC  
**Endpoint:** `tour-service:9004` (RPC port)

#### RPC metode:
- `TourService.CreateTour` - Kreiranje nove ture
- `TourService.GetTours` - Dobijanje tura korisnika

#### API Gateway rute:
- `POST /api/rpc/tours` - RPC kreiranje ture (auth required)
- `GET /api/rpc/tours/my` - RPC dobijanje mojih tura (auth required)

## Struktura fajlova

```
services/
├── auth-service/
│   ├── internal/
│   │   ├── rpc/
│   │   │   └── auth_rpc.go           # RPC service definicija
│   │   └── handler/
│   │       └── auth_rpc_handler.go   # RPC business logic
│   └── cmd/main/main.go              # Pokretanje RPC servera
├── tour-service/
│   ├── internal/
│   │   ├── rpc/
│   │   │   └── tour_rpc.go           # RPC service definicija  
│   │   └── handler/
│   │       └── tour_rpc_handler.go   # RPC business logic
│   └── cmd/main/main.go              # Pokretanje RPC servera
└── api-gateway/
    ├── internal/
    │   ├── rpc/
    │   │   └── rpc_client.go         # RPC klijent
    │   └── handler/
    │       └── rpc_handler.go        # RPC HTTP handler
    └── cmd/main/main.go              # RPC klijent integracija
```

## Testiranje RPC endpoint-a

### Login preko RPC
```bash
curl -X POST http://localhost:8000/api/rpc/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"test123"}'
```

### Kreiranje ture preko RPC
```bash
curl -X POST http://localhost:8000/api/rpc/tours \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Test RPC Tour",
    "description": "Tour kreirana preko RPC protokola",
    "difficulty": "EASY",
    "tags": "rpc,test",
    "distanceKm": 5.0
  }'
```

## Razlike RPC vs REST

| Aspekt | REST API | RPC API |
|--------|----------|---------|
| **Transport** | HTTP/JSON | TCP/Binary |
| **Performanse** | Sporije (HTTP overhead) | Brže (direktni pozivi) |
| **Keširanje** | HTTP keširanje | Nema automatsko keširanje |
| **Debug** | Lako (HTTP tools) | Teže (potrebni specifični tools) |
| **Skalabilnost** | Dobra (stateless) | Dobra (persistent connections) |

## Konfiguracija

### Environment varijable
```
# Auth Service
RPC_PORT=9003

# Tour Service  
RPC_PORT=9004

# API Gateway
AUTH_RPC_ADDR=auth-service:9003
TOUR_RPC_ADDR=tour-service:9004
```

### Docker Compose portovi
```yaml
auth-service:
  ports:
    - "8003:8003"  # HTTP
    - "9003:9003"  # RPC

tour-service:
  ports:
    - "8004:8004"  # HTTP
    - "9004:9004"  # RPC
```

## Implementacijski detalji

- RPC serveri se pokreću paralelno sa HTTP serverima (goroutines)
- Koristi se Go-ov standardni `net/rpc` paket
- RPC handler-i su wrapper-i oko postojeće business logike
- API Gateway ima fallback na HTTP ako RPC nije dostupan
- Greške se propagiraju kroz RPC sa custom error tipovima

Ova implementacija demonstrira hibridni pristup gde sistem podržava i REST i RPC protokole istovremeno.