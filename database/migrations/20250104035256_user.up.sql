CREATE TABLE users (
                       id VARCHAR(26) PRIMARY KEY,
                       email VARCHAR(255) UNIQUE NOT NULL,
                       username VARCHAR(255) NOT NULL,
                       password VARCHAR(255),
                       updated_at TIMESTAMP
)