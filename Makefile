include .env
export

setup:
	docker exec -i app-postgres psql --username $(DB_USERNAME) --password $(DB_PASSWORD) --dbname $(DB_DATABASE) < setup.sql;
up:
	docker-compose up --detach

down:
	docker-compose down --remove-orphans