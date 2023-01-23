PHOO
====
> PHP quick high performance HTTP server, 
> under the hood it uses `PHP-FPM` as a PHP process manager and automates its configurations.

Why?
===
> PHP isn't built for async world, so adopting the community, ecosystem and the mindset 
> to be async isn't an easy task,
> but also I want very simple command to run, and it handles everything without too many configurations files,
> today most of the apps are using environment variables and the well-known `.env` file, so why there isn't a tool
> that you can ask to just run and configure everything from a single `.env` file, I don't want to add a hassle for understanding
> how `PHP-FPM` is working or anything else, all what I need it `$ phoo serve`, that's all!

How?
====
> Basically, `phoo` is a simple static-file as well a fastcgi reverse-proxy, but mainly focuses on `PHP`, not only that,
> but also, you can consider `phoo` a supervisor that manages `PHP-FPM` and its configurations to match today's setup.

Usage?
======
> SOON