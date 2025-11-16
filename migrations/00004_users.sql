-- +goose Up
-- +goose StatementBegin
ALTER table users add if not exists is_active boolean default false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER table users drop if exists is_active;
-- +goose StatementEnd
