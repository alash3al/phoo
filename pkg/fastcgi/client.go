package fastcgi

import (
	"errors"
	"github.com/yookoala/gofast"
	"net/http"
	"strings"
)

type Config struct {
	FastCGIServerURL       string
	DefaultScript          string
	RestrictDotFilesAccess bool
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
	c.handler.ServeHTTP(w, r)
}
