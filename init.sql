CREATE TABLE Users (
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

CREATE TABLE Logs (
    id VARCHAR(255) PRIMARY KEY,
    severity TEXT,
    content TEXT,
    timestamp TIMESTAMP
);