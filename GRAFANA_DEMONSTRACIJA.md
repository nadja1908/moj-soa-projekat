# GRAFANA DEMONSTRACIJA KORAK PO KORAK

## PRISTUP GRAFANI:
- URL: http://localhost:3001
- Username: admin
- Password: admin

## KORAK 1: KREIRAJ NOVI DASHBOARD
1. Klikni "+" u levi panel
2. Klikni "Create Dashboard" 
3. Klikni "Add visualization"

## KORAK 2: DODAJ HOST METRIKE
### Panel 1 - HOST CPU
- Query: `100 - (avg by(instance) (rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100)`
- Title: "HOST - CPU Usage %"
- Panel type: Time series
- Unit: Percent (0-100)
- Objašnjenje: "Ovo pokazuje CPU host mašine"

### Panel 2 - HOST RAM  
- Query: `(1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100`
- Title: "HOST - RAM Usage %"
- Panel type: Time series
- Unit: Percent (0-100)
- Objašnjenje: "Ovo je RAM host mašine"

### Panel 3 - HOST Network
- Query: `rate(node_network_receive_bytes_total[5m])`
- Title: "HOST - Network Traffic"
- Panel type: Time series
- Unit: Bytes/sec
- Objašnjenje: "Mrežni saobraćaj host mašine"

## KORAK 3: DODAJ KONTEJNER METRIKE  
### Panel 4 - KONTEJNER RAM
- Query: `container_memory_working_set_bytes{id="/"} / (1024*1024*1024)`
- Title: "CONTAINERS - Total RAM Usage GB"
- Panel type: Stat
- Unit: GB
- Objašnjenje: "Ukupno RAM svih kontejnera"

### Panel 5 - KONTEJNER Network
- Query: `rate(container_network_receive_bytes_total{id="/"}[5m])`
- Title: "CONTAINERS - Network Traffic" 
- Panel type: Time series
- Unit: Bytes/sec
- Objašnjenje: "Mrežni promet kontejnera"

## KORAK 4: DODAJ APPLICATION METRIKE
### Panel 6 - HTTP Requests
- Query: `http_requests_total{job="follower-service"}`
- Title: "APPLICATION - HTTP Requests"
- Panel type: Stat
- Unit: Short
- Objašnjenje: "HTTP zahtevi od .NET aplikacije"

### Panel 7 - Services Status
- Query: `up`
- Title: "SERVICES - Status"
- Panel type: Table
- Unit: Short
- Objašnjenje: "Status svih servisa"

## RAZLIKA HOST vs KONTEJNER:
- **HOST metrike:** Cela fizička/virtualna mašina
- **KONTEJNER metrike:** Samo Docker aplikacije
- **APPLICATION metrike:** Specifični mikroservisi

## DEMONSTRACIJA PROFESORU:
1. "Evo tri nivoa monitoring-a u Grafani"
2. "Gornji redovi - HOST mašina metrike" 
3. "Srednji redovi - KONTEJNER metrike"
4. "Donji redovi - APPLICATION metrike"
5. "Vidite razliku između host sistema i aplikacija"

---

# TRACING - JAEGER DEMONSTRACIJA

## PRISTUP JAEGER-U:
- URL: http://localhost:16686
- Ne treba login

## AKO NEMA TRACE-OVA (NORMALNO JE!):
**Objasni profesoru:**

*"Jaeger je spreman za tracing, ali trenutno servisi ne šalju trace podatke. U production sistemu, svi mikroservisi bi imali OpenTelemetry konfigurisan. Evo šta tracing omogućava:"*

### ŠTA TRACING RADI:
1. **Distributed Tracing:** Prati zahtev kroz sve mikroservise
2. **Performance:** Vidi gde su bottleneck-ovi  
3. **Dependencies:** Mapira komunikaciju između servisa
4. **Error Tracking:** Locira gde se greške dešavaju

### DEMONSTRACIJA KONCEPTA:
1. **Pokažeš Jaeger UI:** "Ovo je gde bi se trace-ovi prikazali"
2. **Objasniš workflow:** 
   - *"Kada korisnik pošalje zahtev na API Gateway"*
   - *"Gateway pozove Auth Service"* 
   - *"Auth pozove Stakeholders Service"*
   - *"Jaeger bi pokazao celu putanju sa timing-om"*

### DEMONSTRIRAJ UMESTO TRACE-OVA:
```powershell
# Pokaži da servisi komuniciraju (logovi)
docker logs moj-soa-projekat-auth-service-1 --tail 10
docker logs moj-soa-projekat-follower-service-1 --tail 10
```

## PROFESORU OBJASNI:
*"Tracing je treći stub observability-ja. Implementirao sam infrastrukturu (Jaeger), a u production-u bi svaki servis slao trace podatke. Trenutno imam METRICS i LOGGING koji pokrivaju monitoring potrebe."*

## ŠTA POKAZUJEŠ PROFESORU:
- **Trace:** Jedan kompletan zahtev kroz sistem
- **Span:** Deo zahteva u jednom servisu  
- **Timeline:** Koliko dugo je svaki deo trajao
- **Dependencies:** Kako servisi komuniciraju

## PRIMER OBJAŠNJENJA:
*"Kada korisnik pošalje zahtev, vidim da je prošao kroz:*
1. *API Gateway (5ms)*
2. *Auth Service (15ms)* 
3. *Stakeholders Service (25ms)*
*Ukupno vreme: 45ms"*

---

# LOGGING DEMONSTRACIJA

## OPCIJA 1: DOCKER LOGS (Jednostavno)
```bash
# Vidi logove auth servisa
docker logs moj-soa-projekat-auth-service-1 --tail 20

# Vidi logove follower servisa  
docker logs moj-soa-projekat-follower-service-1 --tail 20

# Real-time logovi
docker logs moj-soa-projekat-auth-service-1 -f
```

## OPCIJA 2: CENTRALIZED LOGS (Elasticsearch/Kibana)
- URL: http://localhost:5601 (Kibana)
- **Objašnjenje:** "Svi logovi se centralno čuvaju i mogu se pretraživati"

## ŠTA POKAZUJEŠ PROFESORU:
1. **Distributed Logs:** Logovi iz svih servisa na jednom mestu
2. **Search:** Pretražuj greške, specifične zahteve
3. **Filtering:** Filtriraj po servisu, level (INFO, ERROR)
4. **Timeline:** Vidi kada se događaju problemi

## KOMPLETNA OBSERVABILITY DEMONSTRACIJA:
1. **METRICS** (Prometheus + Grafana) → "Koliko resursa troším"
2. **TRACING** (Jaeger) → "Kako zahtevi putuju kroz sistem" 
3. **LOGGING** (Docker/ELK) → "Šta se tačno događa u kodu"