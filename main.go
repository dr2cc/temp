package main

import (
	"context"
	//"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

type Tags []string

// В этом примере преобразуем данные из таблицы postgresql
// в пользовательский тип - структуру Video
type Video struct {
	Id    string
	Title string
	Tags  Tags
	Views int
}

// приводит сложные типы и структуры к простому типу
type Valuer interface {
	Value() (driver.Value, error)
}

// приводит простой тип к сложным типам и структурам Go
type Scanner interface {
	Scan(src any) error
}

// Value — функция реализующая интерфейс driver.Valuer
func (tags Tags) Value() (driver.Value, error) {
	// преобразуем []string в string
	if len(tags) == 0 {
		return "", nil
	}
	return strings.Join(tags, "|"), nil
}

func (tags *Tags) Scan(value interface{}) error {
	// если `value` равен `nil`, будет возвращён пустой массив
	if value == nil {
		*tags = Tags{}
		return nil
	}

	sv, err := driver.String.ConvertValue(value)
	if err != nil {
		return fmt.Errorf("cannot scan value. %w", err)
	}

	v, ok := sv.(string)
	if !ok {
		return errors.New("cannot scan value. cannot convert value to string")
	}
	*tags = strings.Split(v, "|")

	// удаляем кавычки у тегов
	for i, v := range *tags {
		(*tags)[i] = strings.Trim(v, `"`)
	}
	return nil
}

// получим идентификатор, наименование и теги у роликов
func QueryTagVideos(ctx context.Context, db *sqlx.DB, limit int) ([]Video, error) {
	videos := make([]Video, 0, limit)
	// // самый простой
	// query := `
	// 	SELECT video_id, tags ,views
	// 	FROM videos
	// 	ORDER BY views
	// 	LIMIT $1
	// `
	query := `
		SELECT video_id, MAX(tags) as tags, MAX(views) as views
		FROM videos
		GROUP BY video_id
		ORDER BY MAX(views) DESC
		LIMIT $1
		`

	rows, err := db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var v Video
		// все теги должны автоматически преобразоваться в слайс v.Tags
		err = rows.Scan(&v.Id, &v.Tags, &v.Views)
		if err != nil {
			return nil, err
		}
		videos = append(videos, v)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return videos, nil
}

func main() {
	// db, err := sql.Open("sqlite", "video.db")
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()

	// открываем соединение с БД
	// DataSourceName
	dsn := "postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable"
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	//..

	// получим идентификатор, наименование и теги у роликов
	limit := 5
	list, err := QueryTagVideos(context.Background(), db, limit)
	if err != nil {
		panic(err)
	}
	// для теста проверим, какие строки содержит v.Tags
	// выведем по 4 первых тега
	for _, v := range list {
		length := 4
		if len(v.Tags) < length {
			length = len(v.Tags)
		}
		fmt.Println(strings.Join(v.Tags[:length], " # "))
	}
}
