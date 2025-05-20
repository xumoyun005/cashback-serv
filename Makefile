.PHONY: build run test clean swagger migrate-up migrate-down db-create db-drop setup teardown init

BINARY_NAME=cashback-serv

build:
	go build -o $(BINARY_NAME) cmd/main.go

run: build
	./$(BINARY_NAME)

clean:
	go clean
	rm -f $(BINARY_NAME)

swagger-install:
	go install github.com/swaggo/swag/cmd/swag@latest

swagger-generate: swagger-install
	swag init -g cmd/main.go

init: swagger-generate

db-create:
	@echo "Ma'lumotlar bazasini yaratish..."
	PGPASSWORD=$$(grep DB_PASSWORD .env | cut -d '=' -f2) \
	psql -h $$(grep DB_HOST .env | cut -d '=' -f2) \
	-p $$(grep DB_PORT .env | cut -d '=' -f2) \
	-U $$(grep DB_USER .env | cut -d '=' -f2) \
	-c "CREATE DATABASE $$(grep DB_NAME .env | cut -d '=' -f2);"

db-drop:
	@echo "Ma'lumotlar bazasini o'chirish..."
	PGPASSWORD=$$(grep DB_PASSWORD .env | cut -d '=' -f2) \
	psql -h $$(grep DB_HOST .env | cut -d '=' -f2) \
	-p $$(grep DB_PORT .env | cut -d '=' -f2) \
	-U $$(grep DB_USER .env | cut -d '=' -f2) \
	-c "DROP DATABASE IF EXISTS $$(grep DB_NAME .env | cut -d '=' -f2);"

migrate-up:
	@echo "Migratsiyalarni ishga tushirish..."
	goose -dir migrations postgres "host=$$(grep DB_HOST .env | cut -d '=' -f2) \
	port=$$(grep DB_PORT .env | cut -d '=' -f2) \
	user=$$(grep DB_USER .env | cut -d '=' -f2) \
	password=$$(grep DB_PASSWORD .env | cut -d '=' -f2) \
	dbname=$$(grep DB_NAME .env | cut -d '=' -f2) \
	sslmode=disable" up

migrate-down:
	@echo "Migratsiyalarni orqaga qaytarish..."
	goose -dir migrations postgres "host=$$(grep DB_HOST .env | cut -d '=' -f2) \
	port=$$(grep DB_PORT .env | cut -d '=' -f2) \
	user=$$(grep DB_USER .env | cut -d '=' -f2) \
	password=$$(grep DB_PASSWORD .env | cut -d '=' -f2) \
	dbname=$$(grep DB_NAME .env | cut -d '=' -f2) \
	sslmode=disable" down

setup: db-create migrate-up swagger-generate build

teardown: db-drop clean