-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS "Project" (
	"id" uuid NOT NULL UNIQUE DEFAULT uuid_generate_v4(),
	"owner_id" uuid,
	"name" varchar(255) NOT NULL DEFAULT '',
	"creation_date" date,
	"admin_id" uuid,
	PRIMARY KEY ("id")
);


-- +goose Down
DROP TABLE "Project";
