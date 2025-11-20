package main

import (
	"fmt"
	"log"
	"net/http"

	"rag-training-auth-service/internal/handlers"
	postgres "rag-training-auth-service/internal/storage"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	db, err := postgres.Init()
	if err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Ошибка при закрытии соединения с БД: %v", err)
		}
	}()

	authHandler := &handlers.AuthHandler{DB: db.Conn}
	http.HandleFunc("POST /register", authHandler.Register)
	http.HandleFunc("POST /login", authHandler.Login)
	http.HandleFunc("GET /validate", authHandler.Validate)

	fmt.Println("Starting server at port 8090")
	err = http.ListenAndServe(":8090", nil)
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}
