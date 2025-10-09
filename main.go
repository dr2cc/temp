package main

import (
	"fmt"
	"net/http"
	"time"
)

// helloweb - Snippet for sample hello world webapp (Go)

// Это и есть "простая" функция в качестве обработчика
func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

func twet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "TTT")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /", greet)

	// // http.HandlerFunc— это удобный адаптер,
	// // который позволяет простой функции выполнять "контракт"
	// // http.Handler на обработку HTTP-запросов (чтобы это не значило),
	// // упрощая использование простых функций в качестве обработчиков.
	// http.HandleFunc("/", greet)

	// Это образец HandleFunc в случае использования DefaultServeMux (не использует шаблон CRUD)
	// Будет работать только при http.ListenAndServe("localhost:8080", nil)
	// Оставил для понимания
	http.HandleFunc("/two", twet)

	// // The handler is typically nil, in which case [DefaultServeMux] is used.
	// // Обработчик (второй параметр) по умолчанию равен nil, в этом случае используется [DefaultServeMux].
	// // Его использование не рекомендуется, только в простых, тестовых приложениях.
	// // В рабочих приложениях следует использовать http.NewServeMux или сторонние роутеры
	// http.ListenAndServe("localhost:8080", nil)
	http.ListenAndServe("localhost:8080", mux)
}
