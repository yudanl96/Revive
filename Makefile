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

.PHONY: openmysql migrateup migratedown sqlc test
