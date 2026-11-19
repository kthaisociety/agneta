-- +goose Up
CREATE TABLE roles (
    name text PRIMARY KEY CHECK (name ~ '\S')
);

INSERT INTO roles (name) VALUES
    ('Board advisor'),
    ('Board member'),
    ('Tech lead');

CREATE TABLE reports (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name text NOT NULL CHECK (name ~ '\S'),
    date date NOT NULL,
    role text NOT NULL REFERENCES roles (name),
    submitted_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE records (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    report_id bigint NOT NULL REFERENCES reports (id),
    position integer NOT NULL CHECK (position >= 0),
    question text NOT NULL CHECK (question ~ '\S'),
    answer text NOT NULL CHECK (answer ~ '\S'),
    UNIQUE (report_id, position)
);

-- +goose Down
DROP TABLE records;
DROP TABLE reports;
DROP TABLE roles;
