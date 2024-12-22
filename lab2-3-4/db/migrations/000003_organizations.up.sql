CREATE TABLE organizations
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT ,
    name       TEXT      NOT NULL,
    created_at DATETIME NOT NULL DEFAULT current_timestamp
);