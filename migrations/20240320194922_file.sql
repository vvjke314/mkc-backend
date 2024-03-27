-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS file (
	"id" uuid NOT NULL UNIQUE DEFAULT uuid_generate_v4(),
	"project_id" uuid,
	"filename" varchar(255) NOT NULL DEFAULT '',
	"extension" varchar(255) NOT NULL DEFAULT '',
	"size" int NOT NULL DEFAULT 0,
	"file_path" varchar(255) NOT NULL DEFAULT '',
	"update_datetime" TIMESTAMP,
	PRIMARY KEY ("id"),
	FOREIGN KEY (project_id) REFERENCES "project"(id) ON DELETE CASCADE
);


-- +goose Down
DROP TABLE file;

