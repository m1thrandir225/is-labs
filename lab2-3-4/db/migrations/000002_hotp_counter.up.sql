CREATE TABLE hotp_counters (
    user_id INTEGER PRIMARY KEY,
    counter INTEGER NOT NULL DEFAULT 0,
    last_used_timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
