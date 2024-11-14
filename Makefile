swag:
	swag init -g cmd/mkc/main.go
docker-run:
	docker compose up
nc-build:
	go build -o bin/notechecker cmd/notechecker/notechecker.go
nc-run:
	./bin/notecheckerы
build:
	make swag
	go build -o bin/mkc cmd/mkc/main.go
build-migrate:
	go build -o bin/migrations cmd/migrations/migrations.go
migrate:
# if FLAG != '' will do down migration 
	./bin/migrations $(FLAG)
run:
	./bin/mkc