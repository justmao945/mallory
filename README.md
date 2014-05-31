mallory
=============

Yet another http proxy written in golang, including direct and GAE fetcher


Status
=============
* Support direct fetcher, e.g. http, https...
* Support GAE fetcher, only support http and https with port 443. In this mode we need to deploy fake certificates to forward https requests...
* Simple PAC file server

Installation
=============
```sh
go get github.com/justmao945/mallory/cmd/mallory
```

Direct Engine Usage
=============
```sh
# this is the default mode
mallory
2014/04/12 01:56:33 Listen and serve on 127.0.0.1:18087
2014/04/12 01:56:33 	Engine: direct
```

GAE Engine Usage
=============
```sh
# copy config to the default work dir
mkdir ~/.mallory && cp cfg/* ~/.mallory

# before start the proxy server, we'd better upload the GAE remote application
# for details see https://appengine.google.com
cd mallory/

# put your own app id into app.yaml
vim gae_server/app.yaml

# deploy it with go_appengine, https://developers.google.com/appengine/downloads
goapp deploy gae_server/

# this mode need to use the fake CA, default include mallory.crt and mallory.key
mallory -engine=gae -appspot=your_app_id

# or generate the private key and sign the Root CA by yourself
openssl genrsa -out mallory.key 2048
openssl req -new -x509 -days 365 -key mollory.key -out mallory.crt
```

TODO
=============
* Add test
* Optimize response time
* Add appspot IP resolver
* ....


References
=============
* [goproxy][1]
* [mitmproxy][2]
* [goagent][3]
 

[1]: https://github.com/elazarl/goproxy
[2]: http://mitmproxy.org/
[3]: https://github.com/goagent
