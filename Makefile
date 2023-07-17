#*.env file must be inside src folder
ENV_FILE := src/app.env
include $(ENV_FILE)
export

new_migration:
	migrate create -ext sql -dir src/db/migration/ -seq $(name)

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
	cd src/; mockgen -package mockwk -destination worker/mock/distributor.go desly/worker TaskDistributor

sqlc:
	cd src/; sqlc generate

test:
	cd src/; go test -v -cover -short ./...

server:
	cd src/; go run main.go

tidy:
	cd src/; go mod tidy

proto:
	rm -f src/pb/*.go
	rm -f src/doc/swagger/*.swagger.json
	protoc --proto_path=src/proto --go_out=src/pb --go_opt=paths=source_relative \
    --go-grpc_out=src/pb --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=src/pb --grpc-gateway_opt=paths=source_relative \
		--openapiv2_out=src/doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=deslfy \
    src/proto/*.proto
	statik -src=src/doc/swagger -dest=src/doc -ns=api_docs

redis:
	docker run --name redis -p 6379:6379 -d redis:6.0.20-alpine3.18

.PHONY: migrateup migratedown migrateup1 migratedown1 db_docs sqlc test server tidy proto redis new_migration