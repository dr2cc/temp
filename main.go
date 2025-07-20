package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// Простейший пример, хватит и переменной
var site string = ""

// type HandlerFunc func(ResponseWriter, *Request)
//
// The HandlerFunc type is an adapter to allow the use of ordinary functions as HTTP handlers.
// Тип HandlerFunc — это адаптер, позволяющий использовать обычные функции в качестве HTTP-обработчиков.
// If f is a function with the appropriate signature, HandlerFunc(f) is a [Handler] that calls f.
// Если f — функция с соответствующей сигнатурой, HandlerFunc(f) — это [Handler], вызывающий f.
//
// func (f http.HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request)

// функция обрабатывающая POST запросы к конечной точке "/"
func Transmission() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Запросы для поисковика google:
		//  что такое reader в go
		//  http что такое закрытие потока
		//
		// Играюсь.
		// И такой ридер (то, что можно читать) делаю:
		stringsReader := strings.NewReader("Hello World!")

		//summary
		// Тип Reader можно сделать во многих пакетах (к примеру strings.NewReader и bytes.NewReader)
		//DOGMA// В строку его можно преобразовать только из байтового среза (пока приму как догму)
		// поэтому в начале преобразуем его в []byte
		bytes01, _ := io.ReadAll(stringsReader)
		fmt.Println(string(bytes01))
		bytes02, _ := io.ReadAll(bytes.NewReader(bytes01))
		fmt.Println(string(bytes02))

		// Конец игры

		// читаю r (http запрос)
		//
		// Почему тело http запроса имеет тип ReadCloser?
		//
		// ГЛАВНОЕ - так сделан поток данных в http.
		// А сделан он так вот почему:
		//
		// В контексте HTTP, закрытие потока (stream) относится к процессу
		// прекращения передачи данных между клиентом и сервером, обычно в рамках одного TCP-соединения.
		// Это может происходить по инициативе как клиента, так и сервера, и может быть полным или частичным.
		// Закрытие потоков является важной частью управления соединениями и ресурсами в HTTP.
		// Это позволяет:
		// - избежать утечек ресурсов, предотвратить нежелательное поведение,
		// - и обеспечить более эффективную работу приложений,
		// - в некоторых случаях закрытие потока необходимо для завершения работы приложения.
		bodyToBytes, _ := io.ReadAll(r.Body)

		site = string(bodyToBytes)
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
