-- +migrate Up

CREATE TABLE roles (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(255)
);

CREATE TABLE permissions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    type ENUM('backend', 'frontend') NOT NULL DEFAULT 'backend',
    description VARCHAR(255)
);

CREATE TABLE role_permissions (
    role_id INT NOT NULL,
    permission_id INT NOT NULL,
    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
);

ALTER TABLE users ADD COLUMN role_id INT NULL;
ALTER TABLE users DROP COLUMN role;

ALTER TABLE users
    ADD CONSTRAINT fk_users_roles
        FOREIGN KEY (role_id) REFERENCES roles(id);

INSERT INTO roles (name, description)
VALUES
    ('admin', 'Полный доступ к системе'),
    ('logist', 'Доступ к управлению заказами');