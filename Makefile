run:
	./bin/mkc
build:
	go build -o bin/mkc cmd/mkc/main.go
build-migrate:
	go build -o bin/migrations cmd/migrations/migrations.go
migrate:
# if FLAG != '' will do down migration 
	./bin/migrations $(FLAG)
docker-run:
	docker compose up -d
swag:
	swag init -g cmd/main/mkc/main.go