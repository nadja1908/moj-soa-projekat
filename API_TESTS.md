# Test API pozivi

## 1. Registracija korisnika

```bash
# Registracija turiste
curl -X POST http://localhost:8001/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "marko123",
    "password": "password123",
    "email": "marko@example.com",
    "role": "tourist"
  }'

# Registracija vodiča
curl -X POST http://localhost:8001/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "ana_vodic",
    "password": "vodic2024",
    "email": "ana@guide.com",
    "role": "guide"
  }'
```

## 2. Prijava korisnika

```bash
# Prijava admina
curl -X POST http://localhost:8001/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "secret"
  }'

# Prijava turiste
curl -X POST http://localhost:8001/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "marko123", 
    "password": "password123"
  }'
```

## 3. Admin funkcionalnosti

```bash
# Pregled svih korisnika (potreban admin token)
curl -X GET http://localhost:8001/users \
  -H "Authorization: Bearer ADMIN_TOKEN_HERE"

# Blokiranje korisnika (potreban admin token)
curl -X PUT http://localhost:8001/users/2/block \
  -H "Authorization: Bearer ADMIN_TOKEN_HERE"
```

## 4. Blog funkcionalnosti

```bash
# Kreiranje blog posta (potreban token)
curl -X POST http://localhost:8002/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer USER_TOKEN_HERE" \
  -d '{
    "title": "Putovanje po Srbiji",
    "description": "Najlepša mesta u Srbiji",
    "content": "# Putovanje po Srbiji\n\nSrbija ima puno lepih mesta...\n\n## Novi Sad\nNovi Sad je prelep grad..."
  }'

# Pregled svih blog postova
curl -X GET http://localhost:8002/posts

# Pregled određenog blog posta
curl -X GET http://localhost:8002/posts/1

# Dodavanje komentara (potreban token)
curl -X POST http://localhost:8002/posts/1/comments \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer USER_TOKEN_HERE" \
  -d '{
    "commentText": "Odličan post! Hvala na savetima."
  }'

# Lajkovanje posta (potreban token)
curl -X POST http://localhost:8002/posts/1/like \
  -H "Authorization: Bearer USER_TOKEN_HERE"

# Uklanjanje lajka (potreban token)
curl -X DELETE http://localhost:8002/posts/1/like \
  -H "Authorization: Bearer USER_TOKEN_HERE"
```

## 5. Health check

```bash
# Stakeholders service
curl -X GET http://localhost:8001/health

# Blog service  
curl -X GET http://localhost:8002/health
```