CREATE TABLE access_requests
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id     INTEGER   NOT NULL,
    resource_id INTEGER NOT NULL,
    status      TEXT      NOT NULL,
    reason      TEXT      NOT NULL,
    expires_at  DATETIME  NOT NULL,
    created_at  DATETIME NOT NULL DEFAULT current_timestamp,
    FOREIGN KEY (resource_id) REFERENCES resources (id)
);