BINARY_NAME := rt-auth-service

.PHONY: build run clean docker-build docker-run

## build: Собрать бинарник
build:
	go build -o bin/$(BINARY_NAME) ./cmd/main.go

## run: Запустить сервис
run:
	go run ./cmd/main.go

## clean: Удалить бинарник
clean:
	rm -rf bin/

## docker-build: Собрать Docker-образ
docker-build:
	docker build -t vladkutnyakov/$(BINARY_NAME) .

## docker-push: Запушить Docker-образ
docker-push:
	docker push vladkutnyakov/$(BINARY_NAME):latest

## docker-run: Запустить контейнер (порт 8080 → 8080)
docker-run:
	docker run --rm -d -p 8080:8080 --env-file .env rag-training/$(BINARY_NAME)
