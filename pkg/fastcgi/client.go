package fastcgi

import (
	"errors"
	"github.com/yookoala/gofast"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	FastCGIServerURL       string
	DefaultScript          string
	DocumentRoot           string
	RestrictDotFilesAccess bool
	ServeStaticFiles       bool
	FastCGIParams          map[string]string
}

type Client struct {
	config                 Config
	handler                http.Handler
	fastCGIServerNetwork   string
	fastCGIServerAddr      string
	defaultScriptExtension string
}

func New(config Config) (*Client, error) {
	client := Client{
		config: config,
	}

	if err := client.setDefaultScriptAbsPath(); err != nil {
		return nil, err
	}

	if err := client.setFastCGIServerDetails(); err != nil {
		return nil, err
	}

	client.setFastCGIHandler()

	return &client, nil
}

func (c *Client) setFastCGIServerDetails() error {
	urlParts := strings.Split(c.config.FastCGIServerURL, "://")

	if len(urlParts) != 2 {
		return errors.New("invalid 'FastCGI Server Address' specified")
	}

	c.fastCGIServerNetwork = urlParts[0]
	c.fastCGIServerAddr = urlParts[1]

	return nil
}

func (c *Client) setDefaultScriptAbsPath() error {
	abs, err := filepath.Abs(c.config.DefaultScript)
	if err != nil {
		return err
	}

	c.config.DefaultScript = abs
	c.defaultScriptExtension = strings.ToLower(filepath.Ext(c.config.DefaultScript))

	if _, err := os.Stat(abs); os.IsNotExist(err) {
		return err
	}

	return nil
}

func (c *Client) setFastCGIHandler() {
	sessionHandler := gofast.Chain(
		gofast.MapHeader,
		gofast.BasicParamsMap,
		gofast.MapRemoteHost,
		c.addParams(c.config.FastCGIParams),
	)(gofast.BasicSession)

	c.handler = gofast.NewHandler(
		gofast.NewFileEndpoint(c.config.DefaultScript)(sessionHandler),
		gofast.SimpleClientFactory(
			gofast.SimpleConnFactory(
				c.fastCGIServerNetwork,
				c.fastCGIServerAddr,
			),
		),
	)
}

func (c *Client) addParams(params map[string]string) gofast.Middleware {
	return func(inner gofast.SessionHandler) gofast.SessionHandler {
		return func(client gofast.Client, req *gofast.Request) (*gofast.ResponsePipe, error) {
			req.KeepConn = true

			for k, v := range params {
				req.Params[k] = v
			}

			return inner(client, req)
		}
	}
}

func (c *Client) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !c.config.ServeStaticFiles {
		c.handler.ServeHTTP(w, r)
		return
	}

	filename := filepath.Join(
		c.config.DocumentRoot,
		filepath.Clean(r.URL.Path),
	)

	requestedExtension := strings.ToLower(filepath.Ext(filename))

	// No Dot Files like (.env, .htaccess, ... etc)
	// block any dot file access
	if c.config.RestrictDotFilesAccess && strings.HasPrefix(filepath.Base(filename), ".") {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	fstat, err := os.Stat(filename)

	// Unknown Error
	// there is an error that isn't "FILE NOT FOUND"
	// as we will redirect any not found file to the default server
	if err != nil && !os.IsNotExist(err) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// No Error
	// We assume that it is a static file
	if err == nil && !fstat.IsDir() && requestedExtension != c.defaultScriptExtension {
		http.ServeFile(w, r, filename)
		return
	}

	c.handler.ServeHTTP(w, r)
}
