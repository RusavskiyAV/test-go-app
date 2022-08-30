-- -------------------------------------------------------------
-- TablePlus 4.7.1(428)
--
-- https://tableplus.com/
--
-- Database: test
-- Generation Time: 2022-08-30 06:05:08.7820
-- -------------------------------------------------------------


DROP TABLE IF EXISTS "public"."accounts";
-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS accounts_id_seq;

-- Table Definition
CREATE TABLE "public"."accounts" (
                                     "id" int4 NOT NULL DEFAULT nextval('accounts_id_seq'::regclass),
                                     "user_id" int4 NOT NULL,
                                     "balance" int8 NOT NULL DEFAULT 0 CHECK (balance >= 0),
                                     "created_at" timestamp NOT NULL,
                                     PRIMARY KEY ("id")
);

DROP TABLE IF EXISTS "public"."transactions";
-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS transactions_id_seq;

-- Table Definition
CREATE TABLE "public"."transactions" (
                                         "id" int4 NOT NULL DEFAULT nextval('transactions_id_seq'::regclass),
                                         "account_id" int4 NOT NULL,
                                         "amount" int8 NOT NULL CHECK (amount >= 0),
                                         "reason" varchar(255) NOT NULL,
                                         "participant_account_id" int4,
                                         "created_at" timestamp NOT NULL,
                                         PRIMARY KEY ("id")
);

DROP TABLE IF EXISTS "public"."users";
-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS user_id_seq;

-- Table Definition
CREATE TABLE "public"."users" (
                                  "id" int4 NOT NULL DEFAULT nextval('user_id_seq'::regclass),
                                  "first_name" varchar(255) NOT NULL,
                                  "last_name" varchar(255) NOT NULL,
                                  "created_at" timestamp NOT NULL,
                                  "updated_at" timestamp NOT NULL,
                                  PRIMARY KEY ("id")
);

INSERT INTO "public"."users" ("id", "first_name", "last_name", "created_at", "updated_at") VALUES
                                                                                               (1, 'test_user_name', 'test_user_last_name', '2022-08-30 06:00:00.312532', '2022-08-30 06:00:00.312532'),
                                                                                               (2, 'test_user_name2', 'test_user_last_name2', '2022-08-30 06:00:00.312532', '2022-08-30 06:00:00.312532');

ALTER TABLE "public"."accounts" ADD FOREIGN KEY ("user_id") REFERENCES "public"."users"("id");
ALTER TABLE "public"."transactions" ADD FOREIGN KEY ("account_id") REFERENCES "public"."accounts"("id");
