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
('petar_vodic', '$2a$10$zwx73C.axL2CyPD1R55cVO7gSVq48I9I31E0qnxD8dQuTo474a0Hm', 'petar@example.com', 'guide'),
('milica', '$2a$10$zwx73C.axL2CyPD1R55cVO7gSVq48I9I31E0qnxD8dQuTo474a0Hm', 'milica@example.com', 'tourist'),
('stefan', '$2a$10$zwx73C.axL2CyPD1R55cVO7gSVq48I9I31E0qnxD8dQuTo474a0Hm', 'stefan@example.com', 'tourist'),
('nina_vodic', '$2a$10$zwx73C.axL2CyPD1R55cVO7gSVq48I9I31E0qnxD8dQuTo474a0Hm', 'nina@example.com', 'guide');

INSERT IGNORE INTO profiles (user_id, first_name, last_name, profile_image_url, biography, motto)
VALUES 
(1, 'Ana', 'Anić', 'https://i.pravatar.cc/150?img=1', 'Volim da putujem i otkrivam nove destinacije. Posebno me privlače planine i priroda.', 'Svet je knjiga - ko ne putuje čita samo jednu stranicu.'),
(2, 'Marko', 'Marković', 'https://i.pravatar.cc/150?img=12', 'Strastveni putnik i blogger. Delim svoja iskustva sa letovanja i kulturnih tura.', 'Putuj, uči, rasti.'),
(3, 'Jovana', 'Jovanović', 'https://i.pravatar.cc/150?img=5', 'Avanturista koji voli ekstremne sportove i istraživanje skrivenih destinacija.', 'Život je ili velika avantura ili ništa.'),
(5, 'Petar', 'Petrović', 'https://i.pravatar.cc/150?img=15', 'Iskusan vodič specijalizovan za planinske ture, hiking i pešačenje.', 'Putovanje je jedino što kupujemo, a čini nas bogatijima.'),
(6, 'Milica', 'Milić', 'https://i.pravatar.cc/150?img=9', 'Ljubiteljka umetnosti i kulture. Obožavam da istražujem gradove i njihove muzeje.', 'Kultura je most između naroda.'),
(7, 'Stefan', 'Stefanović', 'https://i.pravatar.cc/150?img=14', 'Digitalni nomad koji putuje svetom i radi remote. Delim savete za rad iz različitih zemalja.', 'Svet je moja kancelarija.'),
(8, 'Nina', 'Nikolić', 'https://i.pravatar.cc/150?img=7', 'Vodilja specijalizovana za wellness i spa ture. Verujem u putovanja koja lečе dušu.', 'Putuj da se sretneš sa sobom.');

INSERT IGNORE INTO users (username, password, email, role)
VALUES ('pera_turista',
        '$2a$10$zwx73C.axL2CyPD1R55cVO7gSVq48I9I31E0qnxD8dQuTo474a0Hm',
        'pera@example.com',
        'tourist');

INSERT IGNORE INTO profiles (user_id, first_name, last_name, profile_image_url, biography, motto)
VALUES (
    9,
    'Pera',
    'Perić',
    'http://localhost:8001/images/pera_turista.png',
    'Zaljubljenik u putovanja, posebno u obilaske istorijskih mesta i muzeja.',
    'Svet je knjiga — ko ne putuje, čita samo jednu stranu.'
);


-- Kreiranje tabele za follower relationships
CREATE TABLE IF NOT EXISTS followers (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    follower_id BIGINT NOT NULL,
    following_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (following_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE KEY unique_follow (follower_id, following_id),
    INDEX idx_follower_id (follower_id),
    INDEX idx_following_id (following_id)
);

-- Dodavanje test follower relationships
-- Ana (1) prati Marka (2), Jovanu (3), Milicu (6)
-- Marko (2) prati Anu (1), Stefana (7)
-- Jovana (3) prati Anu (1), Marka (2), Ninu (8)
-- Milica (6) prati Anu (1), Jovanu (3), Stefana (7)
-- Stefan (7) prati Marka (2), Milicu (6)
-- Nina (8) prati Anu (1), Jovanu (3), Petara (5)
INSERT IGNORE INTO followers (follower_id, following_id) VALUES
(1, 2), -- Ana prati Marka
(1, 3), -- Ana prati Jovanu
(1, 6), -- Ana prati Milicu
(2, 1), -- Marko prati Anu
(2, 7), -- Marko prati Stefana
(3, 1), -- Jovana prati Anu
(3, 2), -- Jovana prati Marka
(3, 8), -- Jovana prati Ninu
(6, 1), -- Milica prati Anu
(6, 3), -- Milica prati Jovanu
(6, 7), -- Milica prati Stefana
(7, 2), -- Stefan prati Marka
(7, 6), -- Stefan prati Milicu
(8, 1), -- Nina prati Anu
(8, 3), -- Nina prati Jovanu
(8, 5); -- Nina prati Petara

-- Kreiranje indeksa za performanse
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_profiles_user_id ON profiles(user_id);