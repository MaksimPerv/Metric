package middleware

import (
	"github.com/MaksimPerv/Metric/internal/compresses"
	"net/http"
	"strings"
)

func GzipMiddleware(h http.Handler) http.Handler {
	//log.Println("aaa")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// по умолчанию устанавливаем оригинальный http.ResponseWriter как тот,
		// который будем передавать следующей функции
		//log.Println(r.Header.Get("Content-Encoding"))
		ow := w
		// проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
		acceptEnc := r.Header.Get("Accept-Encoding")
		supportsGzip := false
		for _, enc := range strings.Split(acceptEnc, ",") {
			if strings.TrimSpace(enc) == "gzip" {
				supportsGzip = true
				break
			}
		}
		if supportsGzip {
			//log.Print("ZBS")
			// оборачиваем оригинальный http.ResponseWriter новым с поддержкой сжатия
			cw := compresses.NewCompressWriter(w)
			// меняем оригинальный http.ResponseWriter на новый
			ow = cw
			// не забываем отправить клиенту все сжатые данные после завершения middleware
			defer cw.Close()
		}
		//log.Print(r.Header.Get("Content-Encoding"))
		// проверяем, что клиент отправил серверу сжатые данные в формате gzip
		contentEnc := r.Header.Get("Content-Encoding")
		isGzipped := false
		for _, enc := range strings.Split(contentEnc, ",") {
			if strings.TrimSpace(enc) == "gzip" {
				isGzipped = true
				break
			}
		}
		if isGzipped {
			// оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
			cr, err := compresses.NewCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// меняем тело запроса на новое
			r.Body = cr
			defer cr.Close()
		}
		// передаём управление хендлеру
		h.ServeHTTP(ow, r)
	})
}
