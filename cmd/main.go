package main

import (
	"fmt"
	"log"
	"net/http"

	"rag-training-auth-service/internal/handlers"
	"rag-training-auth-service/internal/storage"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	db, err := storage.Init()
	if err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Ошибка при закрытии соединения с БД: %v", err)
		}
	}()

	storage.RunMigrations(db)

	authHandler := &handlers.AuthHandler{DB: db.Pool}
	http.HandleFunc("POST /register", authHandler.Register)
	http.HandleFunc("POST /login", authHandler.Login)
	http.HandleFunc("GET /validate", authHandler.Validate)

	fmt.Println("Запуск сервера на порту 8090")
	err = http.ListenAndServe(":8090", nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}
