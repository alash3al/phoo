PHOO
====
> modern php application server, it depends on the bullet-proof `php-fpm` but it controls how it is being run.

Examples
========
> Imagine you have a php application uses modern framework like laravel, symfony ... etc
> that app contains a public directory, and that public directory contains the main bootstrap file that 
> serves the incoming requests named as `index.php`.
```shell
# this all what you need to serve a laravel application!
$ phoo serve -r ./public
⇨ http server started on [::]:8000
```

#### But how about changing the address it is listening on to 0.0.0.0:80?
```shell
# no problem
$ phoo serve -r ./public --http 0.0.0.0:80
⇨ http server started on [::]:80
```

#### Sometimes I want to add custom `pnp.ini` settings, is it easy?
```shell
# is this ok for you? ;)
$ phoo serve -r ./public -i display_errors=Off -i another_key=another_value
⇨ http server started on [::]:8000
```
#### I have a high traffic web-app and I want to increase the php workers
```shell
# just increase the workers count
$ phoo serve -r ./public --workers=20
⇨ http server started on [::]:8000
```

#### Hmmmm, but I want to monitor my app via prometheus metrics, I don't want to do it manually
```shell
# no need to do it yourself, this will enable prometheus metrics at the specified `/metrics` path
$ phoo serve -r ./public --metrics "/metrics"
⇨ http server started on [::]:8000
```

#### Wow!, seems `phoo` has a lot of simple flags/configs, is it documented anywhere?
> just run `phoo serve --help` and enjoy it :), you will find that you can also pass flags via `ENV` vars, and it will automatically read `.env` file in the current working directory.

Requirements
============
- `php-fpm`
- a shell access to run `phoo` :D

Installation
============
- Binary installations could be done via the [releases](https://github.com/alash3al/phoo/releases).
- Docker image is available at [`ghcr.io/alash3al/phoo`](https://github.com/alash3al/phoo/pkgs/container/phoo)
  - you can easily `COPY --from=ghcr.io/alash3al/phoo:2.1.8 /usr/bin/phoo /usr/bin/phoo` to run it into your own custom image!

TODOs
=====
- [ ] Add `.env.example` with comments to describe each var
- 