-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS "Project_access" (
	"id" uuid NOT NULL UNIQUE DEFAULT uuid_generate_v4(),
	"project_id" uuid,
	"customer_id" uuid,
	"customer_access" int NOT NULL DEFAULT 0,
	PRIMARY KEY ("id"),
	FOREIGN KEY (project_id) REFERENCES "Project"(id) ON DELETE CASCADE,
	FOREIGN KEY (customer_id) REFERENCES "Customer"(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE "Project_access";

