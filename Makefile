include .env

# UP

up:
	docker compose --file ./docker-compose.yml --env-file ./.env up -d --build --wait

up-all:
	docker-compose -f ./deployments/compose.yaml up -d --build
	docker compose --file ./docker-compose.yml --env-file ./.env up -d --build --wait

up-musicsnap:
	docker-compose up --build musicsnap-svc

up-b:
	docker-compose up --build musicsnap-svc

up-d:
	docker-compose up --build musicsnap-svc -d

# DOWN

down:
	docker-compose down musicsnap-svc

down-c:
	docker-compose down musicsnap-svc

# Observability

up-obs:
	docker-compose -f ./deployments/compose.yaml up -d --build

down-obs:
	docker-compose -f ./deployments/compose.yaml down

# Migrations

migrate-up:
	migrate -path ./services/musicsnap/migrations -database 'postgres://$(MUSICSNAP_POSTGRES_USER):$(MUSICSNAP_POSTGRES_PASSWORD)@$(MUSICSNAP_POSTGRES_HOST_LOCAL):$(MUSICSNAP_POSTGRES_PORT_EXTERNAL)/$(MUSICSNAP_POSTGRES_NAME)?sslmode=disable' up

migrate-down:
	migrate -path ./services/musicsnap/migrations -database 'postgres://$(MUSICSNAP_POSTGRES_USER):$(MUSICSNAP_POSTGRES_PASSWORD)@$(MUSICSNAP_POSTGRES_HOST_LOCAL):$(MUSICSNAP_POSTGRES_PORT_EXTERNAL)/$(MUSICSNAP_POSTGRES_NAME)?sslmode=disable' down 1

migrate-test-up:
	migrate -path ./services/musicsnap/migrations -database 'postgres://postgres:password@localhost:5444/musicsnap_test?sslmode=disable' up

migrate-test-down:
	migrate -path ./services/musicsnap/migrations -database 'postgres://postgres:password@localhost:5444/musicsnap_test?sslmode=disable' down 1
