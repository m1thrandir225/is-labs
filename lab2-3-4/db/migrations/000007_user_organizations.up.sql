CREATE TABLE user_organizations
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT ,
    user_id    INTEGER   NOT NULL,
    org_id     INTEGER NOT NULL,
    role_id    INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT current_timestamp,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY  (org_id) REFERENCES organizations (id),
    FOREIGN KEY (role_id) REFERENCES roles (id)
);