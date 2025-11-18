-- Tours database initialization script

CREATE DATABASE IF NOT EXISTS tour_db;
USE tour_db;

-- Tours tabela
CREATE TABLE tours (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) DEFAULT 0.0,
    difficulty_level ENUM('EASY', 'MODERATE', 'HARD') NOT NULL DEFAULT 'EASY',
    status ENUM('DRAFT', 'PUBLISHED', 'ARCHIVED') NOT NULL DEFAULT 'DRAFT',
    distance_km DECIMAL(10,2) DEFAULT 0.0,
    author_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    published_at TIMESTAMP NULL,
    archived_at TIMESTAMP NULL,
    
    INDEX idx_status (status),
    INDEX idx_author_id (author_id),
    INDEX idx_created_at (created_at)
);

-- Key points tabela
CREATE TABLE key_points (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    tour_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    latitude DECIMAL(10,8) NOT NULL,
    longitude DECIMAL(11,8) NOT NULL,
    order_index INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (tour_id) REFERENCES tours(id) ON DELETE CASCADE,
    INDEX idx_tour_id (tour_id),
    INDEX idx_order_index (order_index),
    UNIQUE KEY unique_tour_order (tour_id, order_index)
);

-- Transport types enum helper table (za reference)
CREATE TABLE transport_types (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT
);

INSERT INTO transport_types (name, description) VALUES
('WALKING', 'Walking on foot'),
('CYCLING', 'Bicycle transportation'),
('CAR', 'Personal vehicle'),
('BUS', 'Public bus transportation'),
('BOAT', 'Water transportation');

-- Tour durations tabela
CREATE TABLE tour_durations (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    tour_id BIGINT NOT NULL,
    transport_type ENUM('WALKING', 'CYCLING', 'CAR', 'BUS', 'BOAT') NOT NULL,
    duration_minutes INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (tour_id) REFERENCES tours(id) ON DELETE CASCADE,
    INDEX idx_tour_id (tour_id),
    UNIQUE KEY unique_tour_transport (tour_id, transport_type)
);

-- Tour executions tabela (za simulaciju kretanja turista)
CREATE TABLE tour_executions (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    tour_id BIGINT NOT NULL,
    tourist_id BIGINT NOT NULL,
    transport_type ENUM('WALKING', 'CYCLING', 'CAR', 'BUS', 'BOAT') NOT NULL,
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP NULL,
    current_latitude DECIMAL(10,8),
    current_longitude DECIMAL(11,8),
    current_position_index INT DEFAULT 0,
    status ENUM('ACTIVE', 'COMPLETED', 'ABANDONED') NOT NULL DEFAULT 'ACTIVE',
    
    FOREIGN KEY (tour_id) REFERENCES tours(id),
    INDEX idx_tour_id (tour_id),
    INDEX idx_tourist_id (tourist_id),
    INDEX idx_status (status),
    INDEX idx_started_at (started_at)
);

-- Position updates tabela (za čuvanje istorije kretanja)
CREATE TABLE position_updates (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    execution_id BIGINT NOT NULL,
    latitude DECIMAL(10,8) NOT NULL,
    longitude DECIMAL(11,8) NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (execution_id) REFERENCES tour_executions(id) ON DELETE CASCADE,
    INDEX idx_execution_id (execution_id),
    INDEX idx_timestamp (timestamp)
);

-- Sample data za testiranje

-- Test tours
INSERT INTO tours (name, description, difficulty_level, status, author_id, price) VALUES
('Novi Sad City Center', 'Discover the heart of Novi Sad with its historic landmarks and vibrant culture', 'EASY', 'PUBLISHED', 1, 15.50),
('Danube River Walk', 'A scenic walk along the beautiful Danube river', 'MODERATE', 'PUBLISHED', 2, 20.00),
('Petrovaradin Fortress Tour', 'Explore the historic fortress and enjoy panoramic views', 'MODERATE', 'DRAFT', 1, 25.00);

-- Key points for Novi Sad City Center tour
INSERT INTO key_points (tour_id, name, description, latitude, longitude, order_index) VALUES
(1, 'Trg Slobode', 'Main square with the Name of Mary Church', 45.2671, 19.8335, 0),
(1, 'Dunavska Street', 'Historic pedestrian zone', 45.2678, 19.8322, 1),
(1, 'Danube Park', 'Beautiful park by the Danube river', 45.2701, 19.8364, 2),
(1, 'Zmaj Jovina Street', 'Cultural street with galleries and cafes', 45.2665, 19.8345, 3);

-- Key points for Danube River Walk
INSERT INTO key_points (tour_id, name, description, latitude, longitude, order_index) VALUES
(2, 'Quay of Europe', 'Starting point at the riverside promenade', 45.2701, 19.8364, 0),
(2, 'Rainbow Bridge', 'Modern bridge with great views', 45.2721, 19.8401, 1),
(2, 'Fisherman Island', 'Small island perfect for relaxation', 45.2745, 19.8456, 2);

-- Tour durations
INSERT INTO tour_durations (tour_id, transport_type, duration_minutes) VALUES
(1, 'WALKING', 120),
(1, 'CYCLING', 60),
(1, 'CAR', 30),
(2, 'WALKING', 90),
(2, 'CYCLING', 45),
(3, 'WALKING', 180),
(3, 'CAR', 45);

-- Update distances for published tours
UPDATE tours SET distance_km = 3.2 WHERE id = 1;
UPDATE tours SET distance_km = 2.8 WHERE id = 2;