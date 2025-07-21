package main

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// Простейший пример, в качестве хранилища хватит и переменной
var site string = ""

// Так как http.ResponseWriter указан без имени поля, он встраивается в тип gzipWriter,
// который содержит все методы этого интерфейса(!!).
// В противном случае нужно было бы описать методы Header и WriteHeader.
// В примере для gzipWriter достаточно переопределить метод Write([]byte) (int, error)

type gzipWriter struct {
	// !!!Так как http.ResponseWriter указан без имени поля, он встраивается в тип gzipWriter
	http.ResponseWriter
	Writer io.Writer
}

// Вот тут переопределяем метод Write([]byte) (int, error)
// Если его закомментировать, то у gzipWriter появятся все методы http.ResponseWriter
func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}

// обработчик
func defaultHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, "<html><body>"+strings.Repeat("Hello, world<br>", 20)+"</body></html>")
}

func gzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// проверяем, что клиент поддерживает gzip-сжатие
		// это упрощённый пример. В реальном приложении следует проверять все
		// значения r.Header.Values("Accept-Encoding") и разбирать строку
		// на составные части, чтобы избежать неожиданных результатов
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// если gzip не поддерживается, передаём управление
			// дальше без изменений
			next.ServeHTTP(w, r)
			return
		}

		// создаём gzip.Writer (то чему можно писать) поверх текущего w
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		// закрываем поток
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
		next.ServeHTTP(
			gzipWriter{
				ResponseWriter: w,
				Writer:         gz},
			r)
	})
}

func main() {
	// мультиплексор = роутер
	mux := http.NewServeMux()
	// эндпойнт
	mux.HandleFunc("/", defaultHandle)
	// сервер
	http.ListenAndServe(":3000", gzipHandle(mux))
}

// // Из моей тренировки
// //

// // type HandlerFunc func(ResponseWriter, *Request)
// //
// // The HandlerFunc type is an adapter to allow the use of ordinary functions as HTTP handlers.
// // Тип HandlerFunc — это адаптер, позволяющий использовать обычные функции в качестве HTTP-обработчиков.
// // If f is a function with the appropriate signature, HandlerFunc(f) is a [Handler] that calls f.
// // Если f — функция с соответствующей сигнатурой, HandlerFunc(f) — это [Handler], вызывающий f.
// //
// // func (f http.HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request)

// // функция обрабатывающая POST запросы к конечной точке "/"
// func Transmission() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// Играюсь
// 		// И такой ридер (то, что можно читать) делаю
// 		// Два запроса для google:
// 		//  что такое reader в go
// 		//  body в go имеет тип readcloser
// 		nr := strings.NewReader("Hello World!")
// 		br, _ := io.ReadAll(nr)
// 		fmt.Println(string(br))
// 		// Конец игры

// 		// читаю r
// 		tt, _ := io.ReadAll(r.Body)
// 		site = string(tt)
// 		fmt.Fprint(w, site)
// 	}
// }

// // функция обрабатывающая GET запросы к конечной точке "/"
// func Receiving() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// Так пишу ответ в заголовок- происходит редирект
// 		// Вот это пока (17.07.2025) не знаю до конца
// 		w.Header().Set("Location", site)
// 		//w.WriteHeader(http.StatusTemporaryRedirect)
// 	}
// }

// func main() {
// 	// 1. роутер, то, что умеет работать с путями
// 	router := chi.NewRouter()

// 	// 2. обращение к конечной точке - путь и что по этому пути делать
// 	// параметры:
// 	// - паттерн (шаблон)
// 	// - функция HandlerFunc — это адаптер,
// 	//   позволяющий использовать обычные функции в качестве обработчиков HTTP запросов
// 	router.Post("/", Transmission())
// 	router.Get("/", Receiving())

// 	// 3. сервер (служба слушающая порт и "подающая" ответы на порт- ListenAndServe)
// 	// параметры:
// 	// - адрес - хост и порт
// 	// - handler - обработчик, здесь это роутер)
// 	http.ListenAndServe("localhost:8080", router)
// }

// //
