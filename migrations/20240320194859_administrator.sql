-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS "Administrator" (
	"id" uuid NOT NULL UNIQUE DEFAULT uuid_generate_v4(),
	"name" varchar(255) NOT NULL DEFAULT '',
	"email" varchar(255) NOT NULL DEFAULT '',
	"password" varchar(255) NOT NULL DEFAULT '',
	PRIMARY KEY ("id")
);


-- +goose Down
DROP TABLE Administrator;
