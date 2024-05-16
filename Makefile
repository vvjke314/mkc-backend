swag:
	swag init -g cmd/mkc/main.go
run:
	./bin/mkc
build:
	make swag
	go build -o bin/mkc cmd/mkc/main.go
build-migrate:
	go build -o bin/migrations cmd/migrations/migrations.go
migrate:
# if FLAG != '' will do down migration 
	./bin/migrations $(FLAG)
docker-run:
	docker compose up -d
test-repo:
	make migrate FLAG=-d
	make migrate
	make build
	make run
nc-build:
	go build -o bin/notechecker cmd/notechecker/notechecker.go
nc-run:
	./bin/notechecker