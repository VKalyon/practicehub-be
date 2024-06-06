CREATE TABLE metadata (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    mongoid bytea NOT NULL
);