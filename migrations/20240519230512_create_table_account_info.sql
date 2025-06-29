-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "account_info" (
  id VARCHAR PRIMARY KEY NOT NULL,                -- account id, xid
  app_id varchar NOT NULL,                        -- app id
  account_type varchar NOT NULL,                  -- 帳號類型 (user,vendor,admin)
  regis_mode varchar NOT NULL,                    -- NORMAL,GOOGLE..(third party)
  status varchar NOT NULL,                        -- 當前帳號狀態 (ENABLED/BLOCKED)     

  permission jsonb,                               -- 權限 (0:普通)

  password varchar NOT NULL,                      -- user password
  email varchar NOT NULL,                         -- user email (basic account)
  phone varchar NOT NULL,                         -- user phone (optional)

  ev_status bool NOT NULL DEFAULT false,          -- email 驗證狀態
  pv_status bool NOT NULL DEFAULT false,          -- phone 驗證狀態

  created_at int8 NOT NULL,                       -- timestamp ms
  updated_at int8 NOT NULL,                       -- timestamp ms
  deleted_at int8 NULL,                           -- timestamp ms
  login_at int8 NOT NULL,                         -- timestamp ms
  logout_at int8 NOT NULL                         -- timestamp ms
);

--
CREATE TABLE IF NOT EXISTS "user_profile" (
  account_id VARCHAR,                                       -- account id, xid
  user_name varchar,                                        --
  icon BYTEA,                                               -- 
  description varchar,                                      --
  gender varchar,                                           --
  birthday date,                                            --
  job varchar,                                              --
  country varchar,                                          --
  city varchar,                                             --
  district varchar,                                         --
  zip_code varchar,                                         --
  address varchar,                                          --
  shipping_address jsonb,                                   -- 收貨地址

  language varchar NOT NULL,                                --
  email_noti bool NOT NULL DEFAULT false,                   -- email notification
  phone_noti bool NOT NULL DEFAULT false,                   -- phone notification

  created_at int8 NOT NULL,                                 -- timestamp ms
  updated_at int8 NOT NULL,                                 -- timestamp ms
  deleted_at int8 NULL                                      -- timestamp ms
);

CREATE INDEX account_id_idx ON user_profile (account_id);


CREATE TABLE IF NOT EXISTS "verification_code" (
  account_id VARCHAR,                             -- account id, xid
  action varchar,                                 -- ev_login, ev_setpassword, email_verification, pv_login, pv_setpassword, phone_verfication
  code varchar,                                   -- verify code for email or phone
  token varchar,                                  -- verify token for email or phone

  created_at int8 NOT NULL,                       -- timestamp ms
  updated_at int8 NOT NULL,                       -- timestamp ms
  deleted_at int8 NULL                            -- timestamp ms
);

CREATE UNIQUE INDEX account_id_action_idx ON verification_code (account_id, action);
--
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX account_id_action_idx;
DROP TABLE "verification_code";

DROP INDEX account_id_idx;
DROP TABLE "user_profile";

DROP TABLE "account_info";
-- +goose StatementEnd
