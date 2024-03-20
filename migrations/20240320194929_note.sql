-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS "Note" (
	"id" uuid NOT NULL UNIQUE DEFAULT uuid_generate_v4(),
	"project_id" uuid NOT NULL DEFAULT '',
	"title" varchar(255) NOT NULL DEFAULT '',
	"content" varchar(255) NOT NULL DEFAULT '',
	"upload_datetime" date NOT NULL DEFAULT '',
	"deadline" date NOT NULL DEFAULT '',
	PRIMARY KEY ("id")
);


-- +goose Down
DROP TABLE Note;
