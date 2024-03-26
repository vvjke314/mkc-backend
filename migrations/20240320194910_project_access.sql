-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS project_access (
	"id" uuid NOT NULL UNIQUE DEFAULT uuid_generate_v4(),
	"project_id" uuid,
	"customer_id" uuid,
	"customer_access" int NOT NULL DEFAULT 0,
	PRIMARY KEY ("id"),
	FOREIGN KEY (project_id) REFERENCES "project"(id) ON DELETE CASCADE,
	FOREIGN KEY (customer_id) REFERENCES "customer"(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE project_access;

