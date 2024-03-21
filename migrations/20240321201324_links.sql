-- +goose Up
-- +goose StatementBegin
ALTER TABLE "Project" ADD CONSTRAINT Project_fk1 FOREIGN KEY ("owner_id") REFERENCES "Customer"("id");    
ALTER TABLE "Project" ADD CONSTRAINT Project_fk4 FOREIGN KEY ("admin_id") REFERENCES "Administrator"("id");
ALTER TABLE "Project_access" ADD CONSTRAINT Project_access_fk1 FOREIGN KEY ("project_id") REFERENCES "Project"("id");
ALTER TABLE "Project_access" ADD CONSTRAINT Project_access_fk2 FOREIGN KEY ("customer_id") REFERENCES "Customer"("id");
ALTER TABLE "File" ADD CONSTRAINT File_fk1 FOREIGN KEY ("project_id") REFERENCES "Project"("id");
ALTER TABLE "Note" ADD CONSTRAINT Note_fk1 FOREIGN KEY ("project_id") REFERENCES "Project"("id");
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE "Project" DROP CONSTRAINT Project_fk1;
ALTER TABLE "Project" DROP CONSTRAINT Project_fk4;
ALTER TABLE "Project_access" DROP CONSTRAINT Project_access_fk1;
ALTER TABLE "Project_access" DROP CONSTRAINT Project_access_fk2;
ALTER TABLE "File" DROP CONSTRAINT File_fk1;
ALTER TABLE "Note" DROP CONSTRAINT Note_fk1;
-- +goose StatementEnd