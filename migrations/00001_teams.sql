-- +goose Up
-- +goose StatementBegin
CREATE table IF NOT EXISTS teams (
    id serial primary key,
    name text unique not null,
    created_at timestamp default (timezone('utc', now()))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP table if exists teams;
-- +goose StatementEnd
