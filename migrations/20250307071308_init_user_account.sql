-- +goose Up
-- +goose StatementBegin

INSERT INTO account_info
(id, app_id, regis_mode, "permission", "password", email, phone, ev_status, pv_status, created_at, updated_at, deleted_at, login_at, logout_at, account_type, status)
VALUES('cqfqpbfte875u4eecgc0', '', 'NORMAL', '{"order_read": true, "co_marketing": false, "product_read": true, "order_rewrite": true, "product_rewrite": true, "subscribe_email": false, "can_access_cross_account": true}'::jsonb, 'E10ADC3949BA59ABBE56E057F20F883E', 'admin@kiumi.com.tw', '', true, false, 1721740461289, 1741281006262, 0, 0, 0, 1, 1);

INSERT INTO user_profile
(account_id, user_name, icon, description, gender, birthday, job, country, city, address, shipping_address, "language", email_noti, phone_noti, created_at, updated_at, deleted_at, district, zip_code)
VALUES('cqfqpbfte875u4eecgc0', 'Admin', NULL, 'Admin', 'male', '2025-03-07', 'none', 'na', 'na', 'na', '[]'::jsonb, '中文', true, false, 1721740461399, 1740337632646, 0, NULL, NULL);



-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
delete from account_info where id='cqfqpbfte875u4eecgc0';
delete from user_profile where account_id='cqfqpbfte875u4eecgc0';
-- +goose StatementEnd
