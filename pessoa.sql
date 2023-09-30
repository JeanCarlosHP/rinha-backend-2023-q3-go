-- CREATE DATABASE IF NOT EXISTS rinha;

-- USE rinha;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS Pessoas (
  id uuid DEFAULT uuid_generate_v4 (),
  apelido VARCHAR(32) NOT NULL UNIQUE,
  nome VARCHAR(100) NOT NULL,
  nascimento DATE NOT NULL,
  stack VARCHAR(32)[]
);

-- DROP TABLE pessoas;

-- INSERT INTO "pessoas" ("apelido","nome","nascimento","stack") VALUES ('josé','José Roberto','2000-10-01','{"C#","Node","Oracle"}') RETURNING "id"

-- SELECT * FROM "pessoas", (select array_agg(lower(stack)) from unnest("stack") as s) as el WHERE 'Node' = ANY (el);
SELECT * FROM "pessoas" WHERE LOWER(apelido) LIKE '%Node%' OR LOWER(nome) LIKE '%Node%' OR '%Node%' LIKE ANY (stack)

SELECT * FROM "pessoas" WHERE LOWER(apelido) LIKE '%josé%' OR LOWER(nome) LIKE '%node%' OR EXISTS (
  SELECT 1
  FROM unnest(stack) AS elemento
  WHERE LOWER(elemento) = LOWER('py')
);
