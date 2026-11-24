-- +goose Up
CREATE TABLE reports (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name text NOT NULL CHECK (name ~ '\S'),
    date date NOT NULL,
    role text NOT NULL CHECK (role ~ '\S'),
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
