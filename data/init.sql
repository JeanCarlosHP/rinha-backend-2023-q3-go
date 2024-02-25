CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE FUNCTION generate_searchable(_nome VARCHAR, _apelido VARCHAR, _stack VARCHAR(32)[])
    RETURNS VARCHAR AS $$
    BEGIN
    RETURN LOWER(_nome) || LOWER(_apelido) || _stack;
    END;
$$ LANGUAGE plpgsql IMMUTABLE;

CREATE TABLE peoples (
    id uuid PRIMARY KEY NOT NULL,
    nickname VARCHAR(32) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    birth DATE NOT NULL,
    stack VARCHAR(32)[] NULL,
    searchable VARCHAR GENERATED ALWAYS AS (generate_searchable(name, nickname, stack)) STORED
);

CREATE INDEX idx_peoples_searchable ON peoples USING gist (searchable gist_trgm_ops (siglen='64'));
