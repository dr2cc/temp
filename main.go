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
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
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

// ********************************************************************//

type Video struct {
	Id          string    // video_id
	Title       string    // title
	PublishTime time.Time // publish_time
	Tags        []string  // tags
	Views       int       // views
}

func readVideoCSV(ctx context.Context, db *sql.DB, csvFile string) error {
	// открываем csv файл
	file, err := os.Open(csvFile)
	if err != nil {
		return err
	}
	defer file.Close()
	//var videos []Video

	// со множественной вставкой
	videos := make([]Video, 0, 1000)

	// определим индексы нужных полей
	const (
		Id          = 0 // video_id
		Title       = 2 // title
		PublishTime = 5 // publish_time
		Tags        = 6 // tags
		Views       = 7 // views
	)

	// конструируем Reader из пакета encoding/csv
	// он умеет читать строки csv-файла
	r := csv.NewReader(file)
	// пропустим первую строку с именами полей
	if _, err := r.Read(); err != nil {
		return err
	}

	for {
		// csv.Reader за одну операцию Read() считывает одну csv-запись
		l, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		// инициализируем целевую структуру,
		// в которую будем делать разбор csv-записи
		v := Video{
			Id:    l[Id],
			Title: l[Title],
		}
		// парсинг строковых записей в типизированные поля структуры
		if v.PublishTime, err = time.Parse(time.RFC3339, l[PublishTime]); err != nil {
			return err
		}
		tags := strings.Split(l[Tags], "|")
		for i, v := range tags {
			tags[i] = strings.Trim(v, `"`)
		}
		v.Tags = tags
		if v.Views, err = strconv.Atoi(l[Views]); err != nil {
			return err
		}
		// добавляем полученную структуру в слайс
		videos = append(videos, v)

		// со множественной вставкой
		if len(videos) == 1000 {
			if err = insertVideos(ctx, db, videos); err != nil {
				return err
			}
			videos = videos[:0]
		}
	}

	// добавляем оставшиеся записи
	return insertVideos(ctx, db, videos)
}

func insertVideos(ctx context.Context, db *sql.DB, videos []Video) error {
	// начинаем транзакцию
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// можно вызвать Rollback (откат изменений) в defer,
	// если Commit (сохранение изменеий) будет раньше,
	// то Rollback проигнорируется
	defer tx.Rollback()

	// В postgresql  плейсхолдеры (или «заполнители») это $N - $1,$2, ...
	// В sqlite это ?,?, ...

	// PrepareContext - создание (не выполнение!) sql query
	stmt, err := db.PrepareContext(ctx,
		"INSERT INTO videos (video_id, title, publish_time, tags, views)"+
			" VALUES($1,$2,$3,$4,$5)")

	// // db.ExecContext - выполнение запроса (sql query)
	// _, err := db.ExecContext(ctx,
	// 	"INSERT INTO videos (video_id, title, publish_time, tags, views)"+
	// 		" VALUES($1,$2,$3,$4,$5)", v.Id, v.Title, v.PublishTime,
	// 	strings.Join(v.Tags, `|`), v.Views)

	// // Добавить в конец запроса- проверка на дубликаты
	// "ON CONFLICT (video_id) DO UPDATE SET"+
	// "title = EXCLUDED.title,"+
	// "views = EXCLUDED.views"

	if err != nil {
		return err
	}
	defer stmt.Close()

	// а вот теперь выполнение подготовленного sql query
	for _, v := range videos {
		_, err := stmt.ExecContext(ctx, v.Id, v.Title, v.PublishTime,
			strings.Join(v.Tags, `|`), v.Views)

		if err != nil {
			return err
		}
	}
	// завершаем транзакцию
	return tx.Commit()
}

func DbFunc() {
	// открываем соединение с БД
	// DataSourceName
	dsn := "postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	//..
	ctx := context.Background()
	_, err = db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS videos (
	    "video_id" TEXT,
	    "title" TEXT,
	    "publish_time" TEXT,
	    "tags" TEXT,
	    "views" INTEGER
	  )`)

	if err != nil {
		log.Fatal(err)
	}
	// // читаем записи из файла в слайс []Video вспомогательной функцией
	//videos, err := readVideoCSV(".\\USvideos.csv")

	// со множественной вставкой
	err = readVideoCSV(ctx, db, ".\\USvideos.csv")
	if err != nil {
		log.Fatal(err)
	}

	// // теперь "переехала" в readVideoCSV
	// // записываем []Video в базу данных
	// err = insertVideos(ctx, db, videos)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	//fmt.Printf("Всего csv-записей %v\n", len(videos))
}

func main() {
	// Подключение к db и работа с ней
	DbFunc()
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
