postgres:
	docker run --name postg -p 5432:5432 -e POSTGRES_PASSWORD=root -e POSTGRES_USER=root -d postgres:12-alpine

startdb:
	docker start postg && docker start admin

admin:
	docker run --name admin --name admin --network=host -e PGADMIN_DEFAULT_PASSWORD=123456 -e PGADMIN_DEFAULT_EMAIL=gh@mail.com -d dpage/pgadmin4

migrateup:
	migrate -path db/migrations -database "postgresql://root:root@localhost:5432/bankdb?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://root:root@localhost:5432/bankdb?sslmode=disable" -verbose down

test:
	go1.18beta1  test -v -cover ./...

login:
	docker exec -it postg psql -U root -d bankdb   

mock:
	mockgen -package mockdb -destination db/mock/store.go bank/db/sqlc Store    