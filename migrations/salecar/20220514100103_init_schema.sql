-- +goose Up
-- +goose StatementBegin
CREATE TABLE "car_type" (
    "id" smallserial PRIMARY KEY,
    "name" varchar(255) NOT NULL
);

CREATE TABLE "car_mark" (
    "id" serial PRIMARY KEY,
    "name" varchar(255) NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now()),
    "car_type_id" bigint NOT NULL,
    "name_rus" varchar(255) DEFAULT NULL
);

ALTER TABLE "car_mark" ADD FOREIGN KEY ("car_type_id") REFERENCES "car_type" ("id");

CREATE INDEX ON "car_mark" ("car_type_id");

COMMENT ON COLUMN "car_mark"."created_at" IS 'дата создание строки';
COMMENT ON COLUMN "car_mark"."updated_at" IS 'дата обновления строки';

CREATE TABLE "car_model" (
    "id" serial PRIMARY KEY,
    "car_mark_id" bigint NOT NULL,
    "name" varchar(255) NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now()),
    "car_type_id" bigint NOT NULL,
    "name_rus" varchar(255) DEFAULT NULL,
    "is_popular" BOOLEAN NOT NULL DEFAULT FALSE
);

ALTER TABLE "car_model" ADD FOREIGN KEY ("car_type_id") REFERENCES "car_type" ("id");

CREATE INDEX ON "car_model" ("car_type_id");

ALTER TABLE "car_model" ADD FOREIGN KEY ("car_mark_id") REFERENCES "car_mark" ("id");

CREATE INDEX ON "car_model" ("car_mark_id");

COMMENT ON COLUMN "car_model"."created_at" IS 'дата создание строки';
COMMENT ON COLUMN "car_model"."updated_at" IS 'дата обновления строки';

CREATE TABLE "car_generation" (
    "id" bigserial PRIMARY KEY,
    "name" varchar(255) NOT NULL,
    "car_model_id" bigint NOT NULL,
    "begin_year" timestamptz NOT NULL DEFAULT('0001-01-01 00:00:00Z'),
    "end_year" timestamptz NOT NULL DEFAULT('0001-01-01 00:00:00Z'),
    "car_type_id" bigint NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "car_generation" ADD FOREIGN KEY ("car_type_id") REFERENCES "car_type" ("id");

CREATE INDEX ON "car_generation" ("car_type_id");

ALTER TABLE "car_generation" ADD FOREIGN KEY ("car_model_id") REFERENCES "car_model" ("id");

CREATE INDEX ON "car_generation" ("car_model_id");

COMMENT ON COLUMN "car_generation"."created_at" IS 'дата создание строки';
COMMENT ON COLUMN "car_generation"."updated_at" IS 'дата обновления строки';

COMMENT ON COLUMN "car_generation"."begin_year" IS 'начало года';
COMMENT ON COLUMN "car_generation"."end_year" IS 'конец года';

CREATE TABLE "car_serie" (
    "id" bigserial PRIMARY KEY,
    "car_model_id" bigint NOT NULL,
    "name" varchar(255) NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now()),
    "car_generation_id" bigint NOT NULL,
    "car_type_id" bigint NOT NULL
);

ALTER TABLE "car_serie" ADD FOREIGN KEY ("car_type_id") REFERENCES "car_type" ("id");

CREATE INDEX ON "car_serie" ("car_type_id");

ALTER TABLE "car_serie" ADD FOREIGN KEY ("car_model_id") REFERENCES "car_model" ("id");

CREATE INDEX ON "car_serie" ("car_model_id");

ALTER TABLE "car_serie" ADD FOREIGN KEY ("car_generation_id") REFERENCES "car_generation" ("id");

CREATE INDEX ON "car_serie" ("car_generation_id");

COMMENT ON COLUMN "car_serie"."created_at" IS 'дата создание строки';
COMMENT ON COLUMN "car_serie"."updated_at" IS 'дата обновления строки';

CREATE TABLE "car_modification" (
    "id" bigserial PRIMARY KEY,
    "car_serie_id" bigint NOT NULL,
    "car_model_id" bigint NOT NULL,
    "name" varchar(255) NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now()),
    "car_type_id" bigint NOT NULL
);

ALTER TABLE "car_modification" ADD FOREIGN KEY ("car_type_id") REFERENCES "car_type" ("id");

CREATE INDEX ON "car_modification" ("car_type_id");

ALTER TABLE "car_modification" ADD FOREIGN KEY ("car_model_id") REFERENCES "car_model" ("id");

CREATE INDEX ON "car_modification" ("car_model_id");

ALTER TABLE "car_modification" ADD FOREIGN KEY ("car_serie_id") REFERENCES "car_serie" ("id");

CREATE INDEX ON "car_modification" ("car_serie_id");

COMMENT ON COLUMN "car_modification"."created_at" IS 'дата создание строки';
COMMENT ON COLUMN "car_modification"."updated_at" IS 'дата обновления строки';

CREATE TABLE "car_characteristic" (
    "id" bigserial PRIMARY KEY,
    "name" varchar(255) DEFAULT NULL,
    "car_characteristic_id" bigint DEFAULT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now()),
    "car_type_id" bigint NOT NULL
);

ALTER TABLE "car_characteristic" ADD FOREIGN KEY ("car_type_id") REFERENCES "car_type" ("id");

CREATE INDEX ON "car_characteristic" ("car_type_id");

ALTER TABLE "car_characteristic" ADD FOREIGN KEY ("car_characteristic_id") REFERENCES "car_characteristic" ("id");

CREATE INDEX ON "car_characteristic" ("car_characteristic_id");

COMMENT ON COLUMN "car_characteristic"."created_at" IS 'дата создание строки';
COMMENT ON COLUMN "car_characteristic"."updated_at" IS 'дата обновления строки';

CREATE TABLE "car_characteristic_value" (
    "id" bigserial PRIMARY KEY,
    "value" varchar(255) DEFAULT NULL,
    "unit" varchar(255) DEFAULT NULL,
    "car_characteristic_id" bigint NULL,
    "car_modification_id" bigint NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now()),
    "car_type_id" bigint NOT NULL
);

ALTER TABLE "car_characteristic_value" ADD FOREIGN KEY ("car_type_id") REFERENCES "car_type" ("id");

CREATE INDEX ON "car_characteristic_value" ("car_type_id");

ALTER TABLE "car_characteristic_value" ADD FOREIGN KEY ("car_characteristic_id") REFERENCES "car_characteristic" ("id");

CREATE INDEX ON "car_characteristic_value" ("car_characteristic_id");

ALTER TABLE "car_characteristic_value" ADD FOREIGN KEY ("car_modification_id") REFERENCES "car_modification" ("id");

CREATE INDEX ON "car_characteristic_value" ("car_modification_id");

COMMENT ON COLUMN "car_characteristic_value"."created_at" IS 'дата создание строки';
COMMENT ON COLUMN "car_characteristic_value"."updated_at" IS 'дата обновления строки';

CREATE TABLE "car_equipment" (
    "id" bigserial PRIMARY KEY,
    "name" varchar(255) NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now()),
    "car_modification_id" bigint NULL,
    "price_min" numeric DEFAULT NULL,
    "car_type_id" bigint NOT NULL,
    "year" timestamptz NULL DEFAULT('0001-01-01 00:00:00Z')
);

ALTER TABLE "car_equipment" ADD FOREIGN KEY ("car_type_id") REFERENCES "car_type" ("id");

CREATE INDEX ON "car_equipment" ("car_type_id");

ALTER TABLE "car_equipment" ADD FOREIGN KEY ("car_modification_id") REFERENCES "car_modification" ("id");

CREATE INDEX ON "car_equipment" ("car_modification_id");

COMMENT ON COLUMN "car_equipment"."created_at" IS 'дата создание строки';
COMMENT ON COLUMN "car_equipment"."updated_at" IS 'дата обновления строки';

CREATE TABLE "car_option" (
    "id" bigserial PRIMARY KEY,
    "name" varchar(255) NOT NULL,
    "car_option_id" bigint DEFAULT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now()),
    "car_type_id" bigint NOT NULL
);

ALTER TABLE "car_option" ADD FOREIGN KEY ("car_type_id") REFERENCES "car_type" ("id");

CREATE INDEX ON "car_option" ("car_type_id");

ALTER TABLE "car_option" ADD FOREIGN KEY ("car_option_id") REFERENCES "car_option" ("id");

CREATE INDEX ON "car_option" ("car_option_id");

COMMENT ON COLUMN "car_option"."created_at" IS 'дата создание строки';
COMMENT ON COLUMN "car_option"."updated_at" IS 'дата обновления строки';

CREATE TABLE "car_option_value" (
    "id" bigserial PRIMARY KEY,
    "is_base" boolean NOT NULL DEFAULT FALSE,
    "car_option_id" bigint NOT NULL,
    "car_equipment_id" bigint NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now()),
    "car_type_id" bigint NOT NULL
);

ALTER TABLE "car_option_value" ADD FOREIGN KEY ("car_type_id") REFERENCES "car_type" ("id");

CREATE INDEX ON "car_option_value" ("car_type_id");

ALTER TABLE "car_option_value" ADD FOREIGN KEY ("car_option_id") REFERENCES "car_option" ("id");

CREATE INDEX ON "car_option_value" ("car_option_id");

ALTER TABLE "car_option_value" ADD FOREIGN KEY ("car_equipment_id") REFERENCES "car_equipment" ("id");

CREATE INDEX ON "car_option_value" ("car_equipment_id");

COMMENT ON COLUMN "car_option_value"."created_at" IS 'дата создание строки';
COMMENT ON COLUMN "car_option_value"."updated_at" IS 'дата обновления строки';
-- +goose StatementEnd




-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "car_option_value" CASCADE;
DROP TABLE IF EXISTS "car_option";
DROP TABLE IF EXISTS "car_equipment" CASCADE;
DROP TABLE IF EXISTS "car_characteristic_value" CASCADE;
DROP TABLE IF EXISTS "car_characteristic" CASCADE;
DROP TABLE IF EXISTS "car_modification" CASCADE;
DROP TABLE IF EXISTS "car_generation" CASCADE;
DROP TABLE IF EXISTS "car_serie" CASCADE;
DROP TABLE IF EXISTS "car_model" CASCADE;
DROP TABLE IF EXISTS "car_mark" CASCADE;
DROP TABLE IF EXISTS "car_type";
-- +goose StatementEnd
