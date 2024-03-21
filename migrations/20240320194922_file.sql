-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS "File" (
	"id" uuid NOT NULL UNIQUE DEFAULT uuid_generate_v4(),
	"project_id" uuid,
	"filename" varchar(255) NOT NULL DEFAULT '',
	"extension" varchar(255) NOT NULL DEFAULT '',
	"size" int NOT NULL DEFAULT 0,
	"file_path" varchar(255) NOT NULL DEFAULT '',
	"upload_datetime" date,
	PRIMARY KEY ("id"),
	FOREIGN KEY (project_id) REFERENCES "Project"(id) ON DELETE CASCADE
);


-- +goose Down
DROP TABLE "File";

