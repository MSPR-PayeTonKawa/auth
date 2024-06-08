CREATE TABLE users (
    user_id INT PRIMARY KEY,
    password_hash VARCHAR(255)
);

INSERT INTO users (user_id, password_hash) VALUES (1, '$2a$10$N9mkSHDJCpxBdrvEncG17ecspcuQE6ca24c16y6x4HwP9nnvDDA42');