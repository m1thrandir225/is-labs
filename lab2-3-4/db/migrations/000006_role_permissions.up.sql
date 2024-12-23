CREATE TABLE role_permissions
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    role_id     INTEGER NOT NULL,
    resource_id INTEGER NOT NULL,
    can_read    BOOLEAN            NOT NULL DEFAULT false,
    can_write   BOOLEAN            NOT NULL DEFAULT false,
    can_delete  BOOLEAN            NOT NULL DEFAULT false,
    created_at  TIMESTAMP NOT NULL DEFAULT current_timestamp,
    FOREIGN KEY (role_id) REFERENCES roles (id),
    FOREIGN KEY (resource_id) references resources (id)
);