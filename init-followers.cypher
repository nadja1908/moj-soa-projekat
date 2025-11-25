// Kreiranje početnih follower odnosa
// Ana (ID 1) prati Marka (ID 2) i Jovanu (ID 3)
MERGE (ana:User {id: 1})
MERGE (marko:User {id: 2})
MERGE (ana)-[:FOLLOWS]->(marko);

MERGE (ana2:User {id: 1})
MERGE (jovana:User {id: 3})
MERGE (ana2)-[:FOLLOWS]->(jovana);

// Marko (ID 2) prati Anu (ID 1)
MERGE (marko2:User {id: 2})
MERGE (ana3:User {id: 1})
MERGE (marko2)-[:FOLLOWS]->(ana3);

// Jovana (ID 3) prati Marka (ID 2)
MERGE (jovana2:User {id: 3})
MERGE (marko3:User {id: 2})
MERGE (jovana2)-[:FOLLOWS]->(marko3);
