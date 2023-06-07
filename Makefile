#app.env file must be inside src folder
include src/app.env
export

migrateup:
	migrate -path src/db/migration -database "$(DB_SOURCE)" -verbose up

migrateup1:
	migrate -path src/db/migration -database "$(DB_SOURCE)" -verbose up 1

migratedown:
	migrate -path src/db/migration -database "$(DB_SOURCE)" -verbose down

migratedown1:
	migrate -path src/db/migration -database "$(DB_SOURCE)" -verbose down 1

mock:
	cd src/; mockgen -package mockdb -destination db/mock/store.go desly/db/sqlc Store

sqlc:
	cd src/; sqlc generate

test:
	cd src/; go test -v -cover -short ./...

server:
	cd src/; go run main.go

.PHONY: migrateup migratedown migrateup1 migratedown1 sqlc test server