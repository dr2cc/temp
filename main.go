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

// Строка запуска pg в docker
//
//docker run -e POSTGRES_PASSWORD=qwerty -p 5432:5432 -v sprint3:/var/lib/postgresql/data -d postgres

package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	_ "github.com/lib/pq"
	//_ "github.com/mattn/go-sqlite3"
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

// ВСПОМНИТЬ! Про HandlerFunc и (w http.ResponseWriter, r *http.Request)
func HealthCheckHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. подключение

		//DriverName
		dn := "postgres"
		//dn := "sqlite3"

		// DataSourceName
		dsn := "postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable"
		//dsn := "sqlite.db"

		db, err := sql.Open(dn, dsn)
		if err != nil {
			panic(err)
		}

		// делаем запрос
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		// не забываем освободить ресурс
		defer cancel()

		// В принципе только для целей проверки наличия соединения достаточно db.Ping
		// но он тоже получает контест, только не явно.
		// Яндекс советует всюду использовать context
		// И уже все готово для работы с данными!
		err = db.PingContext(ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			//fmt.Fprint(w, "Error connecting to the database:", err)
			return
		}
		//w.WriteHeader(http.StatusOK)
		//fmt.Fprint(w, dn, " - successfully connected to the database!")
		defer db.Close()

		//...

		// // // QueryRowContext выполняет запрос, который, как ожидается, вернет не более одной строки (в нашем случае запрос:
		// // // ВЫБРАТЬ Количество(все) как count ИЗ videos
		// // // должен вернуть только одну строку- количество
		// // row := db.QueryRowContext(ctx,
		// // 	"SELECT COUNT(*) as count FROM videos")

		// // 2. получение данных из таблицы videos (any db - pg or sqlite)
		// row := db.QueryRowContext(ctx,
		// 	"SELECT title, views, channel_title "+
		// 		"FROM videos ORDER BY views DESC LIMIT 1")
		// var (
		// 	title string
		// 	views int
		// 	chati string
		// )

		// // 3. Scan() "переводит" полученные данные в GO-типы
		// // порядок переменных должен соответствовать порядку колонок в запросе
		// err = row.Scan(&title, &views, &chati)
		// if err != nil {
		// 	panic(err)
		// }
		// //fmt.Println(getDesc(ctx, db, "0EbFotkXOiA"))
		// fmt.Printf("%s | %d | %s \r\n", title, views, chati)
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
	// iter10
	// хендлер GET "/ping", который при запросе проверяет соединение с базой данных.
	// При успешной проверке хендлер должен вернуть HTTP-статус 200 OK, при неуспешной — 500 Internal Server Error.
	router.Get("/ping", HealthCheckHandler())

	// 3. сервер (служба слушающая порт и "подающая" ответы на порт- ListenAndServe)
	// параметры:
	// - адрес - хост и порт
	// - handler - обработчик, здесь это роутер)
	http.ListenAndServe("localhost:8080", router)
}

//
