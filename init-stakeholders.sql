-- Kreiranje tabele za korisnike
CREATE TABLE IF NOT EXISTS users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    role ENUM('guide', 'tourist', 'administrator') NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Kreiranje tabele za profile korisnika
CREATE TABLE IF NOT EXISTS profiles (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNIQUE NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    profile_image_url VARCHAR(500),
    biography TEXT,
    motto VARCHAR(200),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Ubacivanje test korisnika
-- Svi korisnici imaju lozinku: "password123" (hashirana sa bcrypt)
INSERT IGNORE INTO users (username, password, email, role) VALUES 
('ana', '$2a$10$zwx73C.axL2CyPD1R55cVO7gSVq48I9I31E0qnxD8dQuTo474a0Hm', 'ana@example.com', 'tourist'),
('marko', '$2a$10$zwx73C.axL2CyPD1R55cVO7gSVq48I9I31E0qnxD8dQuTo474a0Hm', 'marko@example.com', 'tourist'),
('jovana', '$2a$10$zwx73C.axL2CyPD1R55cVO7gSVq48I9I31E0qnxD8dQuTo474a0Hm', 'jovana@example.com', 'tourist'),
('admin', '$2a$10$zwx73C.axL2CyPD1R55cVO7gSVq48I9I31E0qnxD8dQuTo474a0Hm', 'admin@example.com', 'administrator'),
('petar_vodic', '$2a$10$zwx73C.axL2CyPD1R55cVO7gSVq48I9I31E0qnxD8dQuTo474a0Hm', 'petar@example.com', 'guide');

INSERT IGNORE INTO profiles (user_id, first_name, last_name, profile_image_url, biography, motto)
VALUES 
(1, 'Ana', 'Anić', 'https://i.pravatar.cc/150?img=1', 'Volim da putujem i otkrivam nove destinacije. Posebno me privlače planine i priroda.', 'Svet je knjiga - ko ne putuje čita samo jednu stranicu.'),
(2, 'Marko', 'Marković', 'https://i.pravatar.cc/150?img=12', 'Strastveni putnik i blogger. Delim svoja iskustva sa letovanja i kulturnih tura.', 'Putuj, uči, rasti.'),
(3, 'Jovana', 'Jovanović', 'https://i.pravatar.cc/150?img=5', 'Avanturista koji voli ekstremne sportove i istraživanje skrivenih destinacija.', 'Život je ili velika avantura ili ništa.'),
(5, 'Petar', 'Petrović', 'https://i.pravatar.cc/150?img=15', 'Iskusan vodič specijalizovan za planinske ture, hiking i pešačenje.', 'Putovanje je jedino što kupujemo, a čini nas bogatijima.');

INSERT IGNORE INTO users (username, password, email, role)
VALUES ('pera_turista',
        '$2a$10$zwx73C.axL2CyPD1R55cVO7gSVq48I9I31E0qnxD8dQuTo474a0Hm',
        'pera@example.com',
        'tourist');

INSERT IGNORE INTO profiles (user_id, first_name, last_name, profile_image_url, biography, motto)
VALUES (
    3,
    'Pera',
    'Perić',
    'http://localhost:8001/images/pera_turista.png',
    'Zaljubljenik u putovanja, posebno u obilaske istorijskih mesta i muzeja.',
    'Svet je knjiga — ko ne putuje, čita samo jednu stranu.'
);


-- Kreiranje indeksa za performanse
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_profiles_user_id ON profiles(user_id);