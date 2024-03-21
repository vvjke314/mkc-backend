-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS "Customer" (
	"id" uuid NOT NULL UNIQUE DEFAULT uuid_generate_v4(),
	"first_name" varchar(255) NOT NULL DEFAULT '',
	"second_name" varchar(255) NOT NULL DEFAULT '',
	"login" varchar(255) NOT NULL UNIQUE DEFAULT '',
	"password" varchar(255) NOT NULL DEFAULT '',
	"email" varchar(255) NOT NULL UNIQUE DEFAULT '',
	"type" int NOT NULL DEFAULT 0,
	PRIMARY KEY ("id")
);


-- +goose Down
DROP TABLE "Customer";

