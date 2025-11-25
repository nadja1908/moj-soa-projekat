-- Kreiranje tabele za blog postove
CREATE TABLE IF NOT EXISTS blog_posts (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_blog_user_id (user_id),
    INDEX idx_blog_created_at (created_at)
);

-- Kreiranje tabele za komentare na blog
CREATE TABLE IF NOT EXISTS blog_comments (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    blog_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    comment_text TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_comments_blog_id (blog_id),
    INDEX idx_comments_user_id (user_id)
);

-- Kreiranje tabele za lajkove
CREATE TABLE IF NOT EXISTS blog_likes (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    blog_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_user_blog_like (blog_id, user_id),
    INDEX idx_likes_blog_id (blog_id),
    INDEX idx_likes_user_id (user_id)
);


-- Kreiranje tabele za slike blog postova
CREATE TABLE IF NOT EXISTS blog_images (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    blog_id BIGINT NOT NULL,
    image_url VARCHAR(512) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_images_blog_id (blog_id),
    CONSTRAINT fk_images_blog
        FOREIGN KEY (blog_id) REFERENCES blog_posts(id)
        ON DELETE CASCADE
);

-- Insertovanje sample blog postova
INSERT INTO blog_posts (user_id, title, description, content, created_at) VALUES
(1, 'Putovanje po Srbiji', 'Otkrijte lepote naše zemlje', '# Putovanje po Srbiji\n\nSrbija je zemlja sa bogatom istorijom i prelepim prirodnim krajolicima. Od planina do ravnica, od gradova do sela, svaki kutak naše zemlje ima nešto jedinstveno da ponudi.\n\n## Planine Srbije\n\nNaše planine su idealne za zimske sportove i letnje šetnje. Kopaonik, Zlatibor, Tara - sve su to destinacije koje vredi posetiti.\n\n## Gradovi\n\nBeograd, Novi Sad, Niš - svaki grad ima svoju priču i svoje karakteristike.', NOW() - INTERVAL 5 DAY),

(2, 'Najbolje plaže Crne Gore', 'Letovanje na Jadranu', '# Najbolje plaže Crne Gore\n\nCrna Gora je poznata po svojim prelepim plažama duž Jadranskog mora. Od Budve do Ulcinja, svaka plaža ima svoj šarm.\n\n## Budva\n\nBudva je najpoznatije letovalište sa brojnim plažama i bogatim noćnim životom.\n\n## Kotor\n\nKotor nudi kombinaciju istorije i mora, sa starim gradom koji je pod zaštitom UNESCO-a.', NOW() - INTERVAL 4 DAY),

(1, 'Avantura u Nacionalnom parku Tara', 'Priroda u svom najlepšem izdanju', '# Nacionalni park Tara\n\nTara je jedan od najlepših nacionalnih parkova u Srbiji. Poznata je po svojoj netaknutoj prirodi i raznovrsnoj flori i fauni.\n\n## Aktivnosti\n\n- Planinarenje\n- Rafting na Drini\n- Poseta vidikovcima\n- Posmatranje ptica\n\n## Smeštaj\n\nU parku i okolini postoje brojni smeštajni objekti, od hotela do bungalova.', NOW() - INTERVAL 3 DAY),

(3, 'Kulturna ruta kroz Vojvodinu', 'Multikulturalni dragulj Srbije', '# Vojvodina - Kulturni mozaik\n\nVojvodina je poznata po svojoj multietničkoj strukturi i bogatoj kulturnoj baštini.\n\n## Novi Sad\n\nGlavni grad Vojvodine poznat je po Petrovaradinskoj tvrđavi i festivalu Exit.\n\n## Subotica\n\nGrad secesijske arhitekture sa prelepim zgradama i parkovima.\n\n## Gastronomija\n\nVojvodina je raj za ljubitelje dobre hrane - od riblje čorbe do slavonskih kobasica.', NOW() - INTERVAL 2 DAY),

(2, 'Zlatibor - Planina za sve generacije', 'Idealna destinacija za porodični odmor', '# Zlatibor\n\nZlatibor je jedna od najpopularnijih turističkih destinacija u Srbiji, pogodna za sve uzraste.\n\n## Zimski turizam\n\n- Skijanje na Tornik\n- Sankanje\n- Šetnje po snegu\n\n## Letnji turizam\n\n- Planinarenje\n- Vožnja bicikla\n- Poseta Stopića pećini\n\n## Gastronomija\n\nNe propustite da probate domaće sir, kajmak i pršutu!', NOW() - INTERVAL 1 DAY),

(3, 'Đavolja varoš - Prirodno čudo Srbije', 'Misteriozne stene južne Srbije', '# Đavolja varoš\n\nĐavolja varoš je jedinstvena geomorfološka pojava sa 202 kamene figure na području od 67 hektara.\n\n## Legenda\n\nPostoje brojne legende o nastanku ovih figura, najpoznatija govori o kamenom venčanju.\n\n## Poseta\n\nNajbolje vreme za posetu je proleće i jesen. Park je otvoren tokom cele godine.\n\n## Kako stići\n\nNajbliži grad je Kuršumlija, odakle putevi vode direktno do Đavolje varoši.', NOW()),

(1, 'Fruška Gora - Planina vina i manastira', 'Duhovna i vinarska oaza Vojvodine', '# Fruška Gora\n\nFruška Gora je jedinstvena planina u ravnici, poznata po svojim manastirima i vinogradima.\n\n## Manastiri\n\nNa Fruškoj Gori se nalazi 16 pravoslavnih manastira, od kojih su najpoznatiji:\n- Krušedol\n- Hopovo\n- Novo Hopovo\n- Grgeteg\n\n## Vinarije\n\nRegion je poznat po proizvodnji kvalitetnih vina. Poseta vinarijama i degustacije su must-do aktivnosti.', NOW()),

(2, 'Beograd - Grad koji nikad ne spava', 'Vodič kroz prestolnicu Srbije', '# Beograd\n\nBeograd je grad sa više od 7000 godina istorije, poznat po svojoj živoj atmosferi i kulturnoj sceni.\n\n## Istorijske znamenitosti\n\n- Kalemegdan\n- Hram Svetog Save\n- Knez Mihailova ulica\n- Skadarlija\n\n## Noćni život\n\nBeograd je poznat po splavovima i klubovima koji rade do jutra.\n\n## Gastronomija\n\nOd tradicionalnih kafana do modernih restorana - za svakoga ponešto!', NOW()),

-- Novi blogovi za nove korisnike
(6, 'Muzejski tura po Evropi', 'Vodič za ljubitelje umetnosti', '# Muzeji koji morate posetiti\n\nEvropski muzeji čuvaju najvrednije umetnike dela svetske baštine.\n\n## Louvre, Paris\n\nDom Mona Lize i stotina drugih remek dela.\n\n## British Museum, London\n\nIstorijsko blago celokupnog sveta.\n\n## Prado, Madrid\n\nSpanska umetnost kroz vekove.', NOW() - INTERVAL 6 DAY),

(7, 'Remote rad iz Barselone', 'Digitalni nomad vodič', '# Kako raditi iz Barselone\n\nBarselona je savršeno mesto za digitalne nomade.\n\n## Internet i koworking\n\n- Brza internet konekcija\n- Brojni koworking prostori\n- Kafići sa WiFi\n\n## Kvalitet života\n\n- Odličan javni prevoz\n- Plaže u gradu\n- Bogata kulturna scena', NOW() - INTERVAL 5 DAY),

(8, 'Wellness retreat na Bali', 'Putovanje koje leči dušu', '# Bali - Ostrvo duhovnosti\n\nBali nudi savršenu kombinaciju relaksacije i duhovnog rasta.\n\n## Yoga retreats\n\n- Jutarnja yoga na plaži\n- Meditacija u prirodi\n- Zdravi obroci\n\n## Spa tretmani\n\n- Tradicionalni balijski masaža\n- Aromaterapija\n- Detoks programi', NOW() - INTERVAL 4 DAY),

(6, 'Galerije u New York-u', 'Umetnički vodič kroz Big Apple', '# Umetnost u NYC\n\nNew York je svetska prestonica moderne umetnosti.\n\n## MoMA\n\nModerna i savremena umetnost.\n\n## Guggenheim\n\nJedinstvena arhitektura i izložbe.\n\n## Met Museum\n\nOd antike do savremenog doba.', NOW() - INTERVAL 3 DAY),

(7, 'Rad iz Lisabona', 'Portugal za digitalne nomade', '# Lisabon calling\n\nLisabon postaje omiljeno mesto nomada.\n\n## Startup scena\n\n- Rastući tech sektor\n- Networking eventi\n- Podrška za strane preduzetnice\n\n## Životni troškovi\n\n- Pristupačniji od ostalih EU gradova\n- Odlična hrana\n- Bogat noćni život', NOW() - INTERVAL 2 DAY),

(8, 'Ayurveda u Indiji', 'Tradicionalno lečenje', '# Putovanje u Kerala\n\nIndija nudi autentične ayurveda iskustva.\n\n## Ayurveda tretmani\n\n- Panchakarma detoks\n- Herbalnim tretmanima\n- Personalizovaon ishranom\n\n## Meditacija\n\n- Vipassana retreats\n- Ashram život\n- Yoga teacher training', NOW() - INTERVAL 1 DAY);

-- Insertovanje sample blog slika
INSERT INTO blog_images (blog_id, image_url) VALUES
(1, 'https://images.unsplash.com/photo-1469854523086-cc02fe5d8800?w=800'),
(2, 'https://images.unsplash.com/photo-1507525428034-b723cf961d3e?w=800'),
(3, 'https://images.unsplash.com/photo-1501594907352-04cda38ebc29?w=800'),
(4, 'https://images.unsplash.com/photo-1476514525535-07fb3b4ae5f1?w=800'),
(5, 'https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=800'),
(6, 'https://images.unsplash.com/photo-1511593358241-7eea1f3c84e5?w=800'),
(7, 'https://images.unsplash.com/photo-1510312305653-8ed496efae75?w=800'),
(8, 'https://images.unsplash.com/photo-1477959858617-67f85cf4f1df?w=800'),
-- Slike za nove blogove
(9, 'https://images.unsplash.com/photo-1541961017774-22349e4a1262?w=800'),  -- Muzej
(10, 'https://images.unsplash.com/photo-1539650116574-75c0c6d73c6e?w=800'), -- Barselona
(11, 'https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=800'), -- Bali yoga
(12, 'https://images.unsplash.com/photo-1541961017774-22349e4a1262?w=800'), -- NYC galerije
(13, 'https://images.unsplash.com/photo-1555881400-74d7acaacd8b?w=800'), -- Lisabon
(14, 'https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=800'); -- Ayurveda Indija
