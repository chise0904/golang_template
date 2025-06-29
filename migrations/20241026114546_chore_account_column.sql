-- +goose Up
-- +goose StatementBegin
ALTER TABLE account_info DROP account_type;
ALTER TABLE account_info ADD account_type int8 NOT NULL DEFAULT 2; --user

ALTER TABLE account_info DROP status;
ALTER TABLE account_info ADD status int8 NOT NULL DEFAULT 1; --enable
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE account_info DROP account_type;
ALTER TABLE account_info ADD account_type varchar NOT NULL DEFAULT 'user'; --user

ALTER TABLE account_info DROP status;
ALTER TABLE account_info ADD status varchar NOT NULL DEFAULT 'ENABLE';
-- +goose StatementEnd
