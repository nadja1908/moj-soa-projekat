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

-- Ubacivanje test administratora i vodiča
INSERT IGNORE INTO users (username, password, email, role) VALUES 
('admin', '$2a$10$zwx73C.axL2CyPD1R55cVO7gSVq48I9I31E0qnxD8dQuTo474a0Hm', 'admin@example.com', 'administrator'),
('marko_vodic', '$2a$10$zwx73C.axL2CyPD1R55cVO7gSVq48I9I31E0qnxD8dQuTo474a0Hm', 'marko@example.com', 'guide');
-- Lozinka je "password123" - hashirana sa bcrypt

-- Kreiranje indeksa za performanse
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_profiles_user_id ON profiles(user_id);