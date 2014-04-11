http_proxy.go
=============

Yet another http proxy written in golang, including direct and GAE fetcher


Status
=============
* Still under development
* Support direct fetcher, e.g. http, https...
* Support GAE fetcher, only support http and https with port 443. In this mode we need to deploy fake certificates to forward http requests...

TODO
=============
* Optimization response speed
* Add appspot IP resolver
* Add pac server and config
* ....
