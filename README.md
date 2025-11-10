# SOA Projekat - Turistička Aplikacija

## Opis projekta

Mikroservisna aplikacija za turizam sa funkcionalnostima za registraciju korisnika, blog postove, komentare i lajkove. Uključuje kompletnu React frontend aplikaciju.

## Arhitektura

### Backend Servisi:
- **Stakeholders Service** (Port 8001) - Upravljanje korisnicima
- **Blog Service** (Port 8002) - Blog funkcionalnost

### Frontend:
- **React App** (Port 3000) - Korisnički interfejs

### Baze podataka:
- **MySQL** (Port 3307) - Stakeholders baza
- **MySQL** (Port 3308) - Blog baza

## Pokretanje

### Preduslovi:
- Docker i Docker Compose
- Node.js 16+ (za frontend)
- Go 1.21+ (za lokalni development)

### Pokretanje backend servisa sa Docker-om:
```bash
docker-compose up --build
```

### Pokretanje frontend aplikacije:
```bash
cd frontend
npm install
npm start
```

Aplikacija će biti dostupna na:
- **Frontend**: http://localhost:3000
- **Stakeholders API**: http://localhost:8001  
- **Blog API**: http://localhost:8002

### Pokretanje backend servisa lokalno (development):
```bash
# Stakeholders service
cd services/stakeholders-service
go run cmd/main/main.go

# Blog service  
cd services/blog-service
go run cmd/main/main.go
```

## API Endpoints

### Stakeholders Service (http://localhost:8001)

**Public endpoints:**
- `POST /register` - Registracija korisnika
- `POST /login` - Prijava korisnika
- `GET /health` - Health check

**Protected endpoints:**
- `GET /users` - Pregled svih korisnika (admin only)
- `PUT /users/:id/block` - Blokiranje korisnika (admin only)

### Blog Service (http://localhost:8002)

**Public endpoints:**
- `GET /posts` - Pregled svih blog postova
- `GET /posts/:id` - Pregled određenog posta sa komentarima
- `GET /health` - Health check

**Protected endpoints:**
- `POST /posts` - Kreiranje blog posta
- `POST /posts/:id/comments` - Dodavanje komentara
- `POST /posts/:id/like` - Lajkovanje posta
- `DELETE /posts/:id/like` - Uklanjanje lajka

## Funkcionalnosti (Prva KT)

✅ **Funkcionalnost #1**: Registracija korisnika (Vodiči/Turisti)
✅ **Funkcionalnost #2**: Admin pregled korisnika  
✅ **Funkcionalnost #3**: Blokiranje korisnika
✅ **Funkcionalnost #6**: Kreiranje blog postova
✅ **Funkcionalnost #7**: Komentarisanje blog postova
✅ **Funkcionalnost #8**: Lajkovanje blog postova

## Test Korisnici

**Administrator:**
- Username: `admin`
- Password: `password123`

**Ili se registrujte kao:**
- **🗺️ Vodič** - može kreirati blog postove
- **🧳 Turista** - može čitati, komentarisati i označavati postove

## Frontend Funkcionalnosti

### 🔐 Autentifikacija
- Registracija sa odabirom uloge (Vodič/Turista)
- Prijava/Odjava korisnika
- JWT token upravljanje

### 📝 Blog funkcionalnosti
- Pregled blog postova (dostupno svima)
- Kreiranje postova (vodiči i administratori)
- Označavanje postova kao omiljene
- Komentarisanje (u development)

### ⚙️ Admin panel
- Upravljanje korisnicima (blokiranje/odblokiranje)
- Statistike sistema
- Pregled svih korisnika

## Environment Varijable

```env
# Stakeholders Service
PORT=8001
DB_USER=user
DB_PASS=password
DB_HOST=stakeholders-db
DB_NAME=stakeholders_db
JWT_KEY=moj-tajni-kljuc-za-jwt-tokenizaciju-2024

# Blog Service
PORT=8002
DB_USER=user
DB_PASS=password
DB_HOST=blog-db
DB_NAME=blog_db
JWT_KEY=moj-tajni-kljuc-za-jwt-tokenizaciju-2024
```

## Primer korišćenja

### Registracija korisnika:
```bash
curl -X POST http://localhost:8001/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "petar123",
    "password": "mojasifra",
    "email": "petar@example.com", 
    "role": "tourist"
  }'
```

### Prijava:
```bash
curl -X POST http://localhost:8001/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "petar123",
    "password": "mojasifra"
  }'
```

### Kreiranje blog posta:
```bash
curl -X POST http://localhost:8002/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "Moj prvi blog post",
    "description": "Opis mog putovanja",
    "content": "Detaljno opisujem svoje putovanje..."
  }'
```