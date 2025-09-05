// Из моей тренировки
//

// type HandlerFunc func(ResponseWriter, *Request)
//
// The HandlerFunc type is an adapter to allow the use of ordinary functions as HTTP handlers.
// Тип HandlerFunc — это адаптер, позволяющий использовать обычные функции в качестве HTTP-обработчиков.
// If f is a function with the appropriate signature, HandlerFunc(f) is a [Handler] that calls f.
// Если f — функция с соответствующей сигнатурой, HandlerFunc(f) — это [Handler], вызывающий f.
//
// func (f http.HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request)

package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// Простейший пример, хватит и переменной
var site string = ""

// функция обрабатывающая POST запросы к конечной точке "/"
func Transmission() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Играюсь
		// И такой ридер (то, что можно читать) делаю
		// Два запроса для google:
		//  что такое reader в go
		//  body в go имеет тип readcloser
		nr := strings.NewReader("Hello World!")
		br, _ := io.ReadAll(nr)
		fmt.Println(string(br))
		// Конец игры

		// читаю r
		tt, _ := io.ReadAll(r.Body)
		site = string(tt)
		fmt.Fprint(w, site)
	}
}

// функция обрабатывающая GET запросы к конечной точке "/"
func Receiving() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Так пишу ответ в заголовок- происходит редирект
		// Вот это пока (17.07.2025) не знаю до конца
		w.Header().Set("Location", site)
		//w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func main() {
	// 1. роутер, то, что умеет работать с путями
	router := chi.NewRouter()

	// 2. обращение к конечной точке - путь и что по этому пути делать
	// параметры:
	// - паттерн (шаблон)
	// - функция HandlerFunc — это адаптер,
	//   позволяющий использовать обычные функции в качестве обработчиков HTTP запросов
	router.Post("/", Transmission())
	router.Get("/", Receiving())

	// 3. сервер (служба слушающая порт и "подающая" ответы на порт- ListenAndServe)
	// параметры:
	// - адрес - хост и порт
	// - handler - обработчик, здесь это роутер)
	http.ListenAndServe("localhost:8080", router)
}

//
