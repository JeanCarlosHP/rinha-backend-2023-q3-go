CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS Pessoas (
  id uuid DEFAULT uuid_generate_v4 (),
  apelido VARCHAR(32) NOT NULL UNIQUE,
  nome VARCHAR(100) NOT NULL,
  nascimento DATE NOT NULL,
  stack VARCHAR(32)[]
);