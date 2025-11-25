// Kreiranje korisnika u Neo4j bazi
// Prvo kreiramo node-ove za korisnike

CREATE (ana:User {userId: 1, username: 'Ana', email: 'ana@example.com'})
CREATE (marko:User {userId: 2, username: 'Marko', email: 'marko@example.com'})
CREATE (jovana:User {userId: 3, username: 'Jovana', email: 'jovana@example.com'})
CREATE (petar:User {userId: 4, username: 'Petar', email: 'petar@example.com'})
CREATE (nikola:User {userId: 5, username: 'Nikola', email: 'nikola@example.com'})
CREATE (milica:User {userId: 6, username: 'Milica', email: 'milica@example.com'})
CREATE (stefan:User {userId: 7, username: 'Stefan', email: 'stefan@example.com'})
CREATE (nina:User {userId: 8, username: 'Nina', email: 'nina@example.com'});

// Kreiranje FOLLOWS veza između korisnika
// Ana (1) prati Marka (2), Jovanu (3), Milicu (6)
MATCH (ana:User {userId: 1}), (marko:User {userId: 2})
CREATE (ana)-[:FOLLOWS]->(marko);

MATCH (ana:User {userId: 1}), (jovana:User {userId: 3})
CREATE (ana)-[:FOLLOWS]->(jovana);

MATCH (ana:User {userId: 1}), (milica:User {userId: 6})
CREATE (ana)-[:FOLLOWS]->(milica);

// Marko (2) prati Anu (1), Stefana (7)
MATCH (marko:User {userId: 2}), (ana:User {userId: 1})
CREATE (marko)-[:FOLLOWS]->(ana);

MATCH (marko:User {userId: 2}), (stefan:User {userId: 7})
CREATE (marko)-[:FOLLOWS]->(stefan);

// Jovana (3) prati Anu (1), Marka (2), Ninu (8)
MATCH (jovana:User {userId: 3}), (ana:User {userId: 1})
CREATE (jovana)-[:FOLLOWS]->(ana);

MATCH (jovana:User {userId: 3}), (marko:User {userId: 2})
CREATE (jovana)-[:FOLLOWS]->(marko);

MATCH (jovana:User {userId: 3}), (nina:User {userId: 8})
CREATE (jovana)-[:FOLLOWS]->(nina);

// Milica (6) prati Anu (1), Jovanu (3), Stefana (7)
MATCH (milica:User {userId: 6}), (ana:User {userId: 1})
CREATE (milica)-[:FOLLOWS]->(ana);

MATCH (milica:User {userId: 6}), (jovana:User {userId: 3})
CREATE (milica)-[:FOLLOWS]->(jovana);

MATCH (milica:User {userId: 6}), (stefan:User {userId: 7})
CREATE (milica)-[:FOLLOWS]->(stefan);

// Stefan (7) prati Marka (2), Milicu (6)
MATCH (stefan:User {userId: 7}), (marko:User {userId: 2})
CREATE (stefan)-[:FOLLOWS]->(marko);

MATCH (stefan:User {userId: 7}), (milica:User {userId: 6})
CREATE (stefan)-[:FOLLOWS]->(milica);

// Nina (8) prati Anu (1), Jovanu (3), Petara (5)
MATCH (nina:User {userId: 8}), (ana:User {userId: 1})
CREATE (nina)-[:FOLLOWS]->(ana);

MATCH (nina:User {userId: 8}), (jovana:User {userId: 3})
CREATE (nina)-[:FOLLOWS]->(jovana);

MATCH (nina:User {userId: 8}), (petar:User {userId: 5})
CREATE (nina)-[:FOLLOWS]->(petar);