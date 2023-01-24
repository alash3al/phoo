package serve

import (
	"github.com/labstack/gommon/log"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
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

func assetsCacheMiddleware(config *Config, handlerFunc http.HandlerFunc) (http.HandlerFunc, error) {
	memfs := sync.Map{}

	if err := filepath.WalkDir(config.DocumentRoot, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		if filepath.Ext(config.DocumentRoot) == filepath.Ext(path) {
			return nil
		}

		if strings.HasPrefix(filepath.Base(path), ".") {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		memfs.Store(path, data)

		return nil
	}); err != nil {
		return nil, err
	}

	return func(w http.ResponseWriter, r *http.Request) {
		filename := filepath.Join(config.DocumentRoot, r.URL.Path)
		contents, found := memfs.Load(filename)
		if !found {
			handlerFunc(w, r)
			return
		}

		w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(r.URL.Path)))
		w.Write(contents.([]byte))
	}, nil
}
