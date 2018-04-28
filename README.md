HTTP2FCGI
==========
> Quickly serve any `FastCGI` based application with no hassle.

What?
=======
> http2fcgi is a reverse proxy that will convert the standard `http` request to `fcgi` 
request so it can served by i.e `php`, `python` ... etc.

Why?
====
> I wanted a production ready simple and tiny solution to serve some of my `laravel` based projects.

Help?
=====
```bash
➜  http2fcgi http2fcgi -h
Usage of http2fcgi:
  -ext comma separated list
        the fastcgi file extension(s) comma separated list (default "php")
  -fcgi string
        the fcgi backend to connect to, you can pass more fcgi related params as query params (default "unix:///var/run/php/php7.0-fpm.sock")
  -http string
        the http addres to listen on (default ":6065")
  -index string
        the default index file (default "index.php,index.html")
  -listing
        whether to allow directory listing or not
  -root string
        the document root (default "./")
  -router string
        the router filename incase of any 404 error (default "index.php")
  -rtimeout int
        the read timeout, zero means unlimited
  -wtimeout int
        the write timeout, zero means unlimited
```

Download
==========