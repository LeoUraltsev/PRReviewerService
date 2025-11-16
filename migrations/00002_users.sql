-- +goose Up
-- +goose StatementBegin
CREATE TABLE if not exists users(
    id text unique,
    username text not null,
    team_name text not null,
    created_at timestamp default (timezone('utc', now())),
    FOREIGN KEY (team_name) REFERENCES teams (name) ON DELETE CASCADE
);
CREATE INDEX idx_team_name ON users (team_name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP table if exists users;
-- +goose StatementEnd
