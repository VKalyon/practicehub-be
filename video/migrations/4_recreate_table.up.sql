CREATE TABLE IF NOT EXISTS metadata (
     id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    mongoid bytea NOT NULL
);