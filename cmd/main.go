package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	postgres "rag-training-auth-service/internal/storage"

	"golang.org/x/crypto/bcrypt"
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

	ctx := context.Background()

	http.HandleFunc("POST /register", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		var exists bool
		err := db.Conn.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", req.Username).Scan(&exists)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if exists {
			http.Error(w, "User already exists", http.StatusBadRequest)
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = db.Conn.Exec(ctx, "INSERT INTO users (username, password_hash) VALUES ($1, $2)", req.Username, hash)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "User created"})
	})

	fmt.Println("Starting server at port 8090")
	err = http.ListenAndServe(":8090", nil)
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}
