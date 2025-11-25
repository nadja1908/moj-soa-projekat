# Test Unfollow Funkcionalnosti

## Promene napravljene

1. **Neo4j Property Fix** - Promenio `id` u `userId` u sledećim metodama:
   - `FollowUserAsync` 
   - `UnfollowUserAsync`
   
2. **Broj komentara** - Dodao `fetchAllCommentsCount` funkciju koja učitava broj komentara za sve postove odmah nakon učitavanja postova.

## Testiranje

### 1. Provera Neo4j baze
```bash
docker exec -it moj-soa-projekat-neo4j-1 cypher-shell -u neo4j -p neo4jpassword "MATCH (a:User {userId: 1})-[:FOLLOWS]->(b:User) RETURN a.userId, b.userId, b.username;"
```

### 2. Test Unfollow preko Frontend-a
1. Pokreni frontend (ako nije već pokrenut)
2. Uloguj se kao korisnik
3. Idi na Blog Posts stranicu
4. Klikni na "Otprati" dugme pored nekog posta
5. Proveri da li nestaje greška i da li post nestaje sa liste

### 3. Provera broja komentara
1. Osvježi stranicu sa Blog Posts
2. Proveri da li se broj komentara prikazuje odmah (npr "3 komentara")
3. Klikni na "0 komentara" da se proširi - sada bi trebalo da odmah vidiš pravi broj

## Očekivani rezultati

✅ Unfollow radi bez greške
✅ Post nestaje sa liste nakon otpraćivanja
✅ Broj komentara se prikazuje odmah
