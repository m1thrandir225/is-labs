CREATE TABLE roles
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT ,
    name       TEXT      NOT NULL,
    org_id     INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT current_timestamp,
    FOREIGN KEY (org_id) REFERENCES organizations (id)
);