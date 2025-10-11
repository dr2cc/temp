package main

// helloweb - Snippet for sample hello world webapp (Go)
// wr		- Snippet for http Response (Go)

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
)

// Тип реализующий два экземпляра логгера,
// а с методом ServeHTTP он (тип) еще и считается http.Handler
type app struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

func (a *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.infoLogger.Println("I use Handler!")
	fmt.Fprintln(w, "I use Handler!")
}

// Это и есть "простая" функция (the plain function) в качестве обработчика
func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

func newLogger(prefix string) *log.Logger {
	// скопировал у Тузова, не понимаю, что где значит
	// Вроде как os.Stdout это выходной поток (даже толком не знаю, что это)
	//return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	return log.New(os.Stdout, prefix, log.Ldate|log.Ltime)
}

func main() {
	// DefaultServeMux не требует создания экземпляра роутера, только объявление его как nil (http.ListenAndServe("localhost:8080", nil))
	// mux := http.NewServeMux()
	chi := chi.NewRouter()

	example := app{
		infoLogger:  newLogger("INFO: "),
		errorLogger: newLogger("ERROR: "),
	}

	// http.HandlerFunc— это ТИП,
	// удобный адаптер, который позволяет простой функции выполнять "контракт"
	// http.Handler на обработку HTTP-запросов (чтобы это не значило),
	// упрощая использование простых функций в качестве обработчиков.
	//
	// We wrap the plain function `greet` in http.HandlerFunc to make it a Handler
	// Мы оборачиваем простую функцию `greet` в http.HandlerFunc, чтобы сделать ее обработчиком
	gr := http.HandlerFunc(greet)

	// http.Handle("POST /httpHandleFunc", gr)
	chi.Handle("/httpHandleFunc", gr)

	// // HandleFunc это функция которая регистрирует handler для заданного шаблона маршрута
	// http.HandleFunc("/", greet)// образец HandleFunc в случае использования DefaultServeMux
	// mux.HandleFunc("POST /HandleFunc", greet)
	chi.HandleFunc("/HandleFunc", greet)

	// // ❗ В таком виде работает! Что это дает пока не понял..
	// mux.Handle("POST /Handle", &example)
	chi.Handle("/Handle", &example)

	example.infoLogger.Println("The server is starting")

	// // The handler is typically nil, in which case [DefaultServeMux] is used.
	// // Обработчик (второй параметр) по умолчанию равен nil, в этом случае используется [DefaultServeMux].
	// // Его использование не рекомендуется (можно только в простых, тестовых приложениях).
	//
	// // В рабочих приложениях следует использовать http.NewServeMux или сторонние роутеры
	// http.ListenAndServe("localhost:8080", nil)

	// Запуск сервера с обработкой ошибки
	if err := http.ListenAndServe("localhost:8080", chi); err != nil {
		example.errorLogger.Fatal(err)
	}
}
