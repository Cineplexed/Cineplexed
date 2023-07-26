CREATE TABLE users (
    id VARCHAR(255) PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    solved_puzzles INT NOT NULL,
    last_solved_puzzle DATE,
    active BOOLEAN NOT NULL
);

CREATE TABLE logs (
    id VARCHAR(255) PRIMARY KEY,
    severity TEXT,
    content TEXT,
    timestamp TIMESTAMP
);

CREATE TABLE selections (
    date DATE,
    movie TEXT,
    num_correct INT,
    num_incorrect INT,
    tagline TEXT,
    overview TEXT,
    genres TEXT[],
    actors TEXT[],
    revenue INT,
    poster VARCHAR(255),
    year TEXT,
    director TEXT,
    producer TEXT
);