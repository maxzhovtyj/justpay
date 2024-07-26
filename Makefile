up:
	docker-compose up --build -d

migrate:
	psql "postgres://postgres:1111@localhost:5433/" -a -f ./db/schema.sql

