-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS "Note" (
	"id" uuid NOT NULL UNIQUE DEFAULT uuid_generate_v4(),
	"project_id" uuid,
	"title" varchar(255) NOT NULL DEFAULT '',
	"content" varchar(255) NOT NULL DEFAULT '',
	"upload_datetime" date,
	"deadline" date,
	PRIMARY KEY ("id"),
	FOREIGN KEY (project_id) REFERENCES "Project"(id) ON DELETE CASCADE
);


-- +goose Down
DROP TABLE "Note";
