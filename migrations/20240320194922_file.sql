-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS "File" (
	"id" uuid NOT NULL UNIQUE DEFAULT uuid_generate_v4(),
	"project_id" uuid NOT NULL DEFAULT '',
	"filename" varchar(255) NOT NULL DEFAULT '',
	"extension" varchar(255) NOT NULL DEFAULT '',
	"size" int NOT NULL DEFAULT '',
	"file_path" varchar(255) NOT NULL DEFAULT '',
	"upload_datetime" date NOT NULL DEFAULT '',
	PRIMARY KEY ("id")
    ALTER TABLE "Project" ADD CONSTRAINT "Project_fk1" FOREIGN KEY ("owner_id") REFERENCES "Customer"("id");    
    ALTER TABLE "Project" ADD CONSTRAINT "Project_fk4" FOREIGN KEY ("admin_id") REFERENCES "Administrator"("id");
    ALTER TABLE "Project_access" ADD CONSTRAINT "Project_access_fk1" FOREIGN KEY ("project_id") REFERENCES "Project"("id");
    ALTER TABLE "Project_access" ADD CONSTRAINT "Project_access_fk2" FOREIGN KEY ("customer_id") REFERENCES "Customer"("id");
    ALTER TABLE "File" ADD CONSTRAINT "File_fk1" FOREIGN KEY ("project_id") REFERENCES "Project"("id");
    ALTER TABLE "Note" ADD CONSTRAINT "Note_fk1" FOREIGN KEY ("project_id") REFERENCES "Project"("id");
);


-- +goose Down
DROP TABLE "File";

