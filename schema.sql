CREATE TABLE users (
    user_id INT PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255)
);

INSERT INTO users (user_id, email, password_hash) VALUES (1, 'example@example.com', '$2a$10$N9mkSHDJCpxBdrvEncG17ecspcuQE6ca24c16y6x4HwP9nnvDDA42');