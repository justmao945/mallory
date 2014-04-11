http_proxy.go
=============

Yet another http proxy written in golang, including direct and GAE fetcher


Status
=============
* Still under development
* Support direct fetcher, e.g. http, https...
* Support GAE fetcher, only support http and https with port 443. In this mode we need to deploy fake certificates to forward https requests...

Installation
=============
```sh
go get github.com/justmao945/http_proxy.go
```

Direct Engine Usage
=============
```sh
# this is the default mode
http_proxy.go
2014/04/12 01:56:33 Listen and serve on 127.0.0.1:18087
2014/04/12 01:56:33 	Engine: direct
```

GAE Engine Usage
=============
```sh
# before start the proxy server, we'd better upload the GAE remote application
cd http_proxy.go/

# put your own app id into app.yaml
vim gae/app.yaml

# deploy it with go_appengine, https://developers.google.com/appengine/downloads
goapp deploy gae/

# this mode need to use the fake CA, default include mallory.crt and mallory.key
http_proxy.go -engine=gae -appspot=your_app_id -cert=path/to/cert.crt -key=path/to/key.key

# or generate the private key and sign the Root CA by yourself
openssl genrsa -out mallory.key 2048
openssl req -new -x509 -days 365 -key mollory.key -out mallory.crt
```

TODO
=============
* Optimization response speed
* Add appspot IP resolver
* Add pac server and config
* ....
