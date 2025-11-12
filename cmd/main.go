package main

import (
	"fmt"
	"log"
	"net/http"

	postgres "rag-training-auth-service/internal/storage"
)

func main() {
	db, err := postgres.Init()
	if err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Ошибка при закрытии соединения с БД: %v", err)
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, Go!")
	})

	fmt.Println("Starting server at port 8090")
	err = http.ListenAndServe(":8090", nil)
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}
