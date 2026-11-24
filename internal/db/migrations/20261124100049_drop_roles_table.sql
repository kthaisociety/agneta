-- +goose Up
ALTER TABLE reports DROP CONSTRAINT reports_role_fkey;
ALTER TABLE reports ADD CONSTRAINT reports_role_check CHECK (role ~ '\S');
DROP TABLE roles;

-- +goose Down
ALTER TABLE reports DROP CONSTRAINT reports_role_check;

CREATE TABLE roles (
    name text PRIMARY KEY CHECK (name ~ '\S')
);

INSERT INTO roles (name) VALUES
    ('Board advisor'),
    ('Board member'),
    ('Tech lead');

ALTER TABLE reports ADD CONSTRAINT reports_role_fkey FOREIGN KEY (role) REFERENCES roles (name);