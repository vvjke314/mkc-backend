-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS "Project" (
	"id" uuid NOT NULL UNIQUE DEFAULT uuid_generate_v4(),
	"owner_id" uuid,
	"capacity" integer NOT NULL DEFAULT 0,
	"name" varchar(255) NOT NULL DEFAULT '',
	"creation_date" date,
	"admin_id" uuid,
	PRIMARY KEY ("id"),
	FOREIGN KEY (owner_id) REFERENCES "Customer"(id) ON DELETE CASCADE,
	FOREIGN KEY (admin_id) REFERENCES "Administrator"(id) ON DELETE CASCADE
);


-- +goose Down
DROP TABLE "Project";
