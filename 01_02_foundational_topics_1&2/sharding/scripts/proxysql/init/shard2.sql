CREATE TABLE to_do_list (
    id BIGINT PRIMARY KEY,
    user_id INT NOT NULL,
    task VARCHAR(255) NOT NULL
);