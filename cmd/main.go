package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Регистрируем обработчик для всех запросов
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, Go!")
	})

	// Запускаем сервер на порту 8090
	fmt.Println("Starting server at port 8090")
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}
