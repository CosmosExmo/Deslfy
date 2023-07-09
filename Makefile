#*.env file must be inside src folder
ENV_FILE := src/app.env
include $(ENV_FILE)
export

migrateup:
	migrate -path src/db/migration -database "$(DB_SOURCE)" -verbose up

migrateup1:
	migrate -path src/db/migration -database "$(DB_SOURCE)" -verbose up 1

migratedown:
	migrate -path src/db/migration -database "$(DB_SOURCE)" -verbose down

migratedown1:
	migrate -path src/db/migration -database "$(DB_SOURCE)" -verbose down 1

db_docs:
	dbdocs build src/doc/db.dbml

mock:
	cd src/; mockgen -package mockdb -destination db/mock/store.go desly/db/sqlc Store

sqlc:
	cd src/; sqlc generate

test:
	cd src/; go test -v -cover -short ./...

server:
	cd src/; go run main.go

.PHONY: migrateup migratedown migrateup1 migratedown1 db_docs sqlc test server