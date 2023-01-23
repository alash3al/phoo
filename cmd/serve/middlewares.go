package serve

import (
	"github.com/labstack/gommon/log"
	"net/http"
)

func loggerMiddleware(enable bool, handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if enable {
			log.Infoj(map[string]interface{}{
				"host": r.Host,
				"uri":  r.URL.RequestURI(),
			})
		}

		handlerFunc(w, r)
	}
}

func recoverMiddleware(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer (func() {
			if err := recover(); err != nil {
				log.Error(err)
			}
		})()

		handlerFunc(w, r)
	}
}
