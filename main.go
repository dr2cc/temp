package main

// helloweb - Snippet for sample hello world webapp (Go)
// wr		- Snippet for http Response (Go)

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Тип реализующий два экземпляра логгера
// а с методом ServeHTTP он (тип) еще и считается http.Handler
type app struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

func (a *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.infoLogger.Println("I use Handler!")
	fmt.Fprintln(w, "I use Handler!")
}

// На 10.10.25 Теперь сделаю метод для app соответствующий Handler
// И затем сам Handler....

// Это и есть "простая" функция в качестве обработчика
func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

func newLogger() *log.Logger {
	// скопировал у Тузова, не понимаю, что где значит
	// Вроде как os.Stdout это выходной поток (даже толком не знаю, что это)
	return log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	//return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
}

func main() {
	mux := http.NewServeMux()

	example := app{
		infoLogger: newLogger(),
	}

	// http.HandlerFunc— это удобный адаптер,
	// который позволяет простой функции выполнять "контракт"
	// http.Handler на обработку HTTP-запросов (чтобы это не значило),
	// упрощая использование простых функций в качестве обработчиков.
	mux.HandleFunc("POST /HandleFunc", greet)

	// ❗ В таком виде работает! Что еэто дает пока не понял, спешу
	mux.Handle("POST /Handle", &example)

	// // Это образец HandleFunc в случае использования DefaultServeMux (не использует шаблон CRUD)
	// // Будет работать только с http.ListenAndServe("localhost:8080", nil)
	// http.HandleFunc("/", greet)//❗
	// // The handler is typically nil, in which case [DefaultServeMux] is used.
	// // Обработчик (второй параметр) по умолчанию равен nil, в этом случае используется [DefaultServeMux].
	// //
	// // Его использование не рекомендуется, только в простых, тестовых приложениях.
	// // В рабочих приложениях следует использовать http.NewServeMux или сторонние роутеры
	// http.ListenAndServe("localhost:8080", nil)//❗
	//
	example.infoLogger.Println("The server is starting")
	http.ListenAndServe("localhost:8080", mux)
}
