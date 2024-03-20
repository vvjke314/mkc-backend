-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS "Project_access" (
	"id" uuid NOT NULL UNIQUE DEFAULT uuid_generate_v4(),
	"project_id" uuid NOT NULL DEFAULT '',
	"customer_id" uuid NOT NULL DEFAULT '',
	"customer_access" int NOT NULL DEFAULT '',
	PRIMARY KEY ("id")
);

-- +goose Down
DROP TABLE Project_access;

