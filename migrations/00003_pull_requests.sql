-- +goose Up
-- +goose StatementBegin
CREATE Type status_enum as ENUM ('MERGED','OPEN');

CREATE TABLE IF NOT EXISTS pull_requests(
    id text unique not null ,
    name text not null,
    author_id text not null,
    status status_enum not null,
    need_more_reviewers boolean not null,
    created_at timestamp default (timezone('utc', now())),
    merged_at timestamp
);

CREATE TABLE IF NOT EXISTS reviewers (
     pr_id text references pull_requests(id) on delete cascade,
     user_id text not null,
     PRIMARY KEY (pr_id, user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE if exists pull_requests;
DROP TABLE if exists reviewers;
DROP TYPE if exists status_enum;
-- +goose StatementEnd
