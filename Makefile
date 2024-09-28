openmysql:
	docker exec -it revive-mysql mysql -uroot -psecret

migrateup:
	migrate -path db/migration -database "mysql://root:secret@tcp(127.0.0.1:3306)/revive" -verbose up

migratedown:
	migrate -path db/migration -database "mysql://root:secret@tcp(127.0.0.1:3306)/revive" -verbose down

sqlc:
	./sqlc generate

.PHONY: openmysql migrateup migratedown sqlc
