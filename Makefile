SHELL := /bin/bash

.PHONY: all db_clean

db_clean:
	@docker rm shjp_data && echo "Data cleanup complete." || echo "No data container to cleanup"
	@docker stop shjp_db && docker rm shjp_db && echo "DB cleanup complete." || echo "No DB container. Nothing to do."

db_init: db_clean
	@docker create -v /shjp_data --name shjp_data postgres
	@docker run --name shjp_db --volumes-from shjp_data -e POSTGRES_USER=shjp -e POSTGRES_PASSWORD=hellochurch -e POSTGRES_DB=shjp -p 5432:5432 -d postgres
	@echo "DB container initialized"

db_up_dev:
	@goose -dir migrations postgres "user=shjp password=hellochurch host=localhost port=5432 dbname=shjp sslmode=disable" up

db_down_dev:
	@goose -dir migrations postgres "user=shjp password=hellochurch host=localhost port=5432 dbname=shjp sslmode=disable" down

db_up_dev_win:
	@~/go/bin/goose.exe -dir migrations postgres "user=shjp password=hellochurch host=`docker-machine.exe ip` port=5432 dbname=shjp sslmode=disable" up

db_down_dev_win:
	@~/go/bin/goose.exe -dir migrations postgres "user=shjp password=hellochurch host=`docker-machine.exe ip` port=5432 dbname=shjp sslmode=disable" down

db: db_clean db_init db_up_dev

db_win: db_clean db_init db_up_dev_win

db_fixtures:
	@go run cmd/fixtures/main.go

db_fixtures_win:
	@go run cmd/fixtures/main.go --host=`docker-machine.exe ip`

db_reset:
	@goose -dir migrations postgres "user=shjp password=hellochurch host=localhost port=5432 dbname=shjp sslmode=disable" down
	@goose -dir migrations postgres "user=shjp password=hellochurch host=localhost port=5432 dbname=shjp sslmode=disable" up
	@make db_fixtures

server:
	@go run cmd/server/main.go

server_win:
	@go run cmd/server/main.go --host=`docker-machine.exe ip`

db_remote_up:
	@source env.sh
	@goose -dir migrations postgres "user=${SHJP_DB_USER} password=${SHJP_DB_PASSWORD} host=${SHJP_DB_HOST} port=${SHJP_DB_PORT} dbname=${SHJP_DB_DATABASE}" up

db_remote_down:
	@source env.sh
	@goose -dir migrations postgres "user=${SHJP_DB_USER} password=${SHJP_DB_PASSWORD} host=${SHJP_DB_HOST} port=${SHJP_DB_PORT} dbname=${SHJP_DB_DATABASE}" down

local:
	./env.sh && cd cmd/server && go run main.go
