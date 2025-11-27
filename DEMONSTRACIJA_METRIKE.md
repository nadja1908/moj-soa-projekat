# DEMONSTRACIJA METRIKA ZA PROFESORA - FINALNE VERZIJE

## KORAK 1: HOST MAŠINA METRIKE ✅
1. **CPU Iskorišćenje:**
   ```promql
   100 - (avg by(instance) (rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100)
   ```
   *"Procenat CPU koji koristi host mašina"*

2. **RAM Memorija:**
   ```promql
   (1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100
   ```
   *"Procenat iskorišćene RAM memorije na host mašini"*

3. **File System:**
   ```promql
   100 - ((node_filesystem_avail_bytes{mountpoint="/"} / node_filesystem_size_bytes{mountpoint="/"}) * 100)
   ```
   *"Procenat popunjenosti disk prostora"*

4. **Mrežni Saobraćaj:**
   ```promql
   rate(node_network_receive_bytes_total[5m])
   ```
   *"Brzina dolaznog mrežnog saobraćaja u bytes/sekundi"*

## KORAK 2: KONTEJNER METRIKE ✅
1. **Ukupna RAM Potrošnja Kontejnera:**
   ```promql
   container_memory_working_set_bytes{id="/"} / (1024*1024*1024)
   ```
   *"Ukupna RAM potrošnja svih Docker kontejnera u GB (WSL2 nivo)"*

2. **Network Promet Kontejnera:**
   ```promql
   rate(container_network_receive_bytes_total{id="/"}[5m])
   ```
   *"Brzina dolaznog mrežnog prometa kontejnera"*

3. **Disk I/O Kontejnera:**
   ```promql
   rate(container_fs_writes_bytes_total{id="/",device=~"overlay.*"}[5m])
   ```
   *"Brzina pisanja na disk za kontejnere"*

4. **CPU Kontejnera:**
   ```promql
   rate(container_cpu_usage_seconds_total{id="/"}[5m]) * 100
   ```
   *"CPU iskorišćenje kontejnera u procentima"*

## KORAK 3: APPLICATION METRIKE ✅
1. **HTTP Zahtevi (.NET Follower Service):**
   ```promql
   http_requests_total{job="follower-service"}
   ```
   *"Broj HTTP zahteva na follower servis"*

2. **Servisi Status:**
   ```promql
   up
   ```
   *"Status svih servisa (1=radi, 0=ne radi)"*

## OBJAŠNJENJE PROFESORU:
**"Implementirao sam tri nivoa monitoring-a:"**

1. **HOST MAŠINA** (Node Exporter) - CPU, RAM, Disk, Network host sistema
2. **KONTEJNER NIVO** (cAdvisor) - Agregiranu potrošnju svih Docker kontejnera 
3. **APLIKACIJA** (Prometheus.NET) - HTTP metrike iz mikroservisa

**"Na Windows-u sa Docker Desktop-om, kontejneri se vide kao jedan WSL2 kontejner, što je normalno. Važno je da vidimo potrošnju resursa."**

**"Ovo pokriva sve zahteve zadatka za 3 poena!"**