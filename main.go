package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alash3al/go-fastcgi-client"
)

var (
	// FlagHTTPAddr .
	FlagHTTPAddr = flag.String("http", ":6065", "the http address to listen on")

	// FlagDocRoot .
	FlagDocRoot = flag.String("root", "./", "the document root")

	// FlagFCGIBackend .
	FlagFCGIBackend = flag.String("fcgi", "unix:///var/run/php/php7.0-fpm.sock", "the fcgi backend to connect to, you can pass more fcgi related params as query params")

	// FlagIndex .
	FlagIndex = flag.String("index", "index.php,index.html", "the default index file `comma separated list`")

	// FlagRouter .
	FlagRouter = flag.String("router", "index.php", "the router filename incase of any 404 error")

	// FlagAllowListing .
	FlagAllowListing = flag.Bool("listing", false, "whether to allow directory listing or not")

	// FlagFCGIExt .
	FlagFCGIExt = flag.String("ext", "php", "the fastcgi file extension(s) `comma separated list`")

	// FlagReadTimeout .
	FlagReadTimeout = flag.Int("rtimeout", 0, "the read timeout, zero means unlimited")

	// FlagWriteTimeout .
	FlagWriteTimeout = flag.Int("wtimeout", 0, "the write timeout, zero means unlimited")
)

var (
	// FCGIBackendConfig .
	FCGIBackendConfig *BackendConfig
)

// BackendConfig the backend configurations i.e 'ext' or any other fcgi params
type BackendConfig struct {
	Network string
	Address string
	Ext     []string
	Params  map[string]string
}

func main() {
	flag.Parse()

	fmt.Println("⇨ checking the fcgi backend ...")
	cnf, err := GetBackendConfig(*FlagFCGIBackend)

	if err != nil {
		log.Fatal(err)
	}

	FCGIBackendConfig = cnf

	fmt.Printf("⇨ http server started on %s\n", *FlagHTTPAddr)
	log.Fatal(http.ListenAndServe(*FlagHTTPAddr, http.HandlerFunc(Serve)))
}

// Serve the main http handler
func Serve(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	defer func() {
		if err := recover(); err != nil {
			res.WriteHeader(500)
			res.Write([]byte("ForwarderError, please see the logs"))
			log.Println(err)
		}
	}()

	filename := filepath.Join(*FlagDocRoot, req.URL.Path)
	scriptname := req.URL.Path
	isFCGI := false
	tryIndex := false
	dir := *FlagDocRoot

	if !IsValidFile(filename) {
		tryIndex = true
	} else if IsValidDir(filename) {
		tryIndex, dir = true, filename
	} else if IsValidFCGIExt(filepath.Ext(filename)) {
		isFCGI = true
	}

	if tryIndex {
		for _, v := range strings.Split(*FlagIndex, ",") {
			if f := filepath.Join(dir, v); IsValidFile(f) {
				filename = filepath.Join(dir, v)
				scriptname = filepath.Join("/", v)
				if IsValidFCGIExt(filepath.Ext(filename)) {
					isFCGI = true
				}
				break
			} else {
				filename = filepath.Join(dir, *FlagRouter)
				scriptname = filepath.Join("/", *FlagRouter)
				isFCGI = true
			}
		}
	}

	fullfilename, _ := filepath.Abs(filename)
	if fullfilename == "" {
		fullfilename = filename
	}

	if !isFCGI && IsValidDir(fullfilename) && !*FlagAllowListing {
		http.Error(res, "DirectoryListing isn't allowed!", 403)
		return
	}

	if !isFCGI {
		http.ServeFile(res, req, fullfilename)
		return
	}

	if !IsValidFile(fullfilename) {
		res.WriteHeader(404)
		res.Write([]byte("Cannot find the requested resource :("))
		return
	}

	pathInfo := req.URL.Path
	if strings.Contains(pathInfo, *FlagRouter) {
		parts := strings.Split(pathInfo, *FlagRouter)
		if len(parts) < 2 {
			parts = append(parts, "/")
		}
		pathInfo = filepath.Join("/", parts[1])
	}

	host, port, _ := net.SplitHostPort(req.RemoteAddr)
	params := map[string]string{
		"SERVER_SOFTWARE":    "http2fcgi",
		"SERVER_PROTOCOL":    req.Proto,
		"REQUEST_METHOD":     req.Method,
		"REQUEST_TIME":       fmt.Sprintf("%d", time.Now().Unix()),
		"REQUEST_TIME_FLOAT": fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Microsecond)),
		"QUERY_STRING":       req.URL.RawQuery,
		"DOCUMENT_ROOT":      fullfilename,
		"REMOTE_ADDR":        host,
		"REMOTE_PORT":        port,
		"SCRIPT_FILENAME":    fullfilename,
		"PATH_TRANSLATED":    fullfilename,
		"SCRIPT_NAME":        scriptname,
		"REQUEST_URI":        req.URL.RequestURI(),
		"AUTH_DIGEST":        req.Header.Get("Authorization"),
		"PATH_INFO":          pathInfo,
		"ORIG_PATH_INFO":     pathInfo,
	}

	for k, v := range req.Header {
		if len(v) < 1 {
			continue
		}
		k = strings.ToUpper(fmt.Sprintf("HTTP_%s", strings.Replace(k, "-", "_", -1)))
		params[k] = strings.Join(v, ";")
	}

	c, e := fcgiclient.Dial(FCGIBackendConfig.Network, FCGIBackendConfig.Address)
	if c == nil {
		res.WriteHeader(500)
		res.Write([]byte(e.Error()))
		return
	}
	defer c.Close()

	c.SetReadTimeout(time.Duration(*FlagReadTimeout) * time.Second)
	c.SetSendTimeout(time.Duration(*FlagWriteTimeout) * time.Second)

	resp, err := c.Request(params, req.Body)
	if resp == nil || resp.Body == nil || err != nil {
		res.WriteHeader(500)
		res.Write([]byte(err.Error()))
		return
	}
	defer resp.Body.Close()

	for k, vals := range resp.Header {
		for _, v := range vals {
			res.Header().Add(k, v)
		}
	}

	res.Header().Set("Server", "http2fcgi")

	if res.Header().Get("X-SendFile") != "" && IsValidFile(res.Header().Get("X-SendFile")) {
		sendFilename := res.Header().Get("X-SendFile")
		res.Header().Del("X-SendFile")
		http.ServeFile(res, req, sendFilename)
		return
	}

	if resp.ContentLength > 0 {
		res.Header().Set("Content-Length", fmt.Sprintf("%d", resp.ContentLength))
	}

	res.WriteHeader(resp.StatusCode)

	n, _ := io.Copy(res, resp.Body)
	if n < 1 {
		stderr := c.Stderr()
		stderr.WriteTo(res)
	}
}

// GetBackendConfig returns the configs of the fcgi backend
func GetBackendConfig(backend string) (cnf *BackendConfig, err error) {
	var u *url.URL

	u, err = url.Parse(backend)
	if err != nil {
		return nil, err
	}

	cnf = &BackendConfig{}
	cnf.Params = map[string]string{}
	cnf.Ext = []string{}

	if ext := strings.ToLower(*FlagFCGIExt); ext != "" {
		cnf.Ext = strings.Split(ext, ",")
	} else {
		return nil, errors.New("You should specifiy the fastcgi script extension i,e '?ext=php'")
	}

	u.Scheme = strings.ToLower(u.Scheme)

	if u.Scheme == "" && u.Host == "" && u.Path != "" {
		cnf.Network, cnf.Address = "unix", u.Path
	}
	if u.Scheme == "" && u.Host != "" && u.Path == "" {
		cnf.Network, cnf.Address = "tcp", u.Host
	}
	if u.Scheme == "unix" && u.Path != "" {
		cnf.Network, cnf.Address = "unix", u.Path
	}
	if u.Scheme == "tcp" && u.Host != "" {
		cnf.Network, cnf.Address = "tcp", u.Host
	}

	for k, v := range u.Query() {
		if len(v) < 1 {
			v = []string{""}
		}
		cnf.Params[k] = v[0]
	}

	if cnf.Network == "" || cnf.Address == "" {
		return nil, errors.New("Invalid fastcgi address (" + backend + ") specified `malformed`")
	}

	if cnf.Network == "unix" && !IsValidFile(cnf.Address) {
		return nil, errors.New("Invalid fastcgi address (" + backend + ") specified `invalid filename`")
	}

	if cnf.Network == "tcp" && !IsValidHost(cnf.Address) {
		return nil, errors.New("Invalid fastcgi address (" + backend + ") specified `invalid host`")
	}

	return cnf, nil
}

// IsValidFile whether the specified filename is valid
func IsValidFile(filename string) bool {
	if filename == "" {
		return false
	}
	_, err := os.Stat(filename)
	return err == nil
}

// IsValidDir whther the specified directory is valid or not
func IsValidDir(filename string) bool {
	stat, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return stat.IsDir()
}

// IsValidHost whether the specified host is online or not
func IsValidHost(host string) bool {
	if host == "" {
		return false
	}
	timeout := time.Duration(2 * time.Second)
	conn, err := net.DialTimeout("tcp", host, timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// IsValidFCGIExt wehther the specified extension is valid fcgi or not
func IsValidFCGIExt(ext string) bool {
	for _, x := range FCGIBackendConfig.Ext {
		if ext == "."+x {
			return true
		}
	}
	return false
}
