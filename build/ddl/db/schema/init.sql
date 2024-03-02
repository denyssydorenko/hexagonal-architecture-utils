DROP DATABASE IF EXISTS hexagonal_architecture_utils;    

CREATE DATABASE hexagonal_architecture_utils WITH OWNER postgres;  

\c hexagonal_architecture_utils;        

CREATE TABLE IF NOT EXISTS users (
    id uuid DEFAULT gen_random_uuid(), 
    name TEXT NOT NULL,
    surname TEXT NOT NULL,
    country TEXT,
    age INT NOT NULL,
    PRIMARY KEY (id)
); 

INSERT INTO users (name, surname, country, age) VALUES ('Olga', 'Dominguez', 'Spain', 20);
INSERT INTO users (name, surname, age) VALUES ('David', 'Gomez', 20);