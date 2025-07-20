package main

import (
	"io"
	"net/http"
	"strings"
)

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

func defaultHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, "<html><body>"+strings.Repeat("Hello, world<br>", 20)+"</body></html>")
}

func main() {

}
