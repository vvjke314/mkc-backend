-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS note (
	"id" uuid NOT NULL UNIQUE DEFAULT uuid_generate_v4(),
	"project_id" uuid,
	"title" varchar(255) NOT NULL DEFAULT '',
	"content" varchar NOT NULL DEFAULT '',
	"update_datetime" TIMESTAMP,
	"deadline" TIMESTAMP,
	PRIMARY KEY ("id"),
	FOREIGN KEY (project_id) REFERENCES "project"(id) ON DELETE CASCADE
);


-- +goose Down
DROP TABLE note;
