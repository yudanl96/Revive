openmysql:
	docker exec -it revive-mysql mysql -uroot -psecret

migrateup:
	migrate -path db/migration -database "mysql://root:secret@tcp(127.0.0.1:3306)/revive?parseTime=true" -verbose up

migratedown:
	migrate -path db/migration -database "mysql://root:secret@tcp(127.0.0.1:3306)/revive?parseTime=true" -verbose down

sqlc:
	./sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockDB -destination db/mock/store.go github.com/yudanl96/revive/db/sqlc Store

proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out pb \
    --grpc-gateway_opt paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=revive\
    proto/*.proto

evans:
	evans --host localhost --port 50051 -r repl

.PHONY: openmysql migrateup migratedown sqlc test server mock proto evans
