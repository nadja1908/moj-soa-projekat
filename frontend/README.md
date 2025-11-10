# SOA Blog Platform - Frontend

React frontend aplikacija za SOA Blog Platform mikroservisnu arhitekturu.

## 📋 Funkcionalnosti

### 🔐 Autentifikacija
- **Registracija korisnika** sa odabirom uloge (Vodič/Turista)
- **Prijava/Odjava** korisnika
- **JWT token** upravljanje

### 👥 Uloge korisnika
- **🧳 Turista**: Čitanje postova, komentarisanje, označavanje omiljenih
- **🗺️ Vodič**: Sve što i turista + kreiranje blog postova
- **👑 Administrator**: Sve funkcionalnosti + upravljanje korisnicima

### 📝 Blog funkcionalnosti
- **Pregled blog postova** (dostupno svima)
- **Kreiranje postova** (vodiči i administratori)
- **Označavanje postova** kao omiljene (registrovani korisnici)
- **Komentarisanje** (registrovani korisnici)

### ⚙️ Admin panel
- **Upravljanje korisnicima** (blokiranje/odblokiranje)
- **Statistike sistema**
- **Pregled svih korisnika**

## 🚀 Pokretanje

### Preduslovi
- Node.js (v16 ili noviji)
- npm ili yarn
- Pokrenuti backend servisi (localhost:8001 i localhost:8002)

### Instalacija
```bash
cd frontend
npm install
```

### Pokretanje u development modu
```bash
npm start
```

Aplikacija će biti dostupna na: http://localhost:3000

### Build za produkciju
```bash
npm run build
```

## 🏗️ Struktura projekta

```
frontend/
├── public/
│   └── index.html
├── src/
│   ├── components/          # React komponente
│   │   ├── NavigationBar.js
│   │   ├── Login.js
│   │   ├── Register.js
│   │   ├── Dashboard.js
│   │   ├── BlogPosts.js
│   │   ├── CreatePost.js
│   │   └── AdminPanel.js
│   ├── context/             # React Context
│   │   └── AuthContext.js
│   ├── services/            # API servisi
│   │   └── api.js
│   ├── App.js
│   ├── index.js
│   └── index.css
└── package.json
```

## 🔧 Konfiguracija

### API Endpoints
- **Stakeholders Service**: `http://localhost:8001`
- **Blog Service**: `http://localhost:8002`

### Routing
- `/` - Početna stranica (redirect na dashboard ili login)
- `/login` - Prijava
- `/register` - Registracija
- `/dashboard` - Korisnički dashboard
- `/posts` - Pregled blog postova
- `/create-post` - Kreiranje novog posta
- `/admin` - Admin panel (samo administratori)

## 👤 Test nalozi

- **Administrator**: 
  - Username: `admin`
  - Password: `password123`

- **Registracija**: Možete kreirati nove naloge sa ulogom Vodič ili Turista

## 🎨 UI/UX
- **Bootstrap 5** za styling
- **React Bootstrap** komponente
- **Responsive design**
- **Intuitivna navigacija**
- **Role-based interface**

## 🔒 Sigurnost
- JWT token autentifikacija
- Automatska odjava na 401 greške
- Role-based access control
- Validacija na frontend i backend nivou

## 🌐 Browser podrška
- Chrome (latest)
- Firefox (latest)
- Safari (latest)
- Edge (latest)

## 📱 Responzivnost
Aplikacija je potpuno responzivna i prilagođena svim veličinama ekrana.