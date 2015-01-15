mallory
=============

[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/justmao945/mallory?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Yet another http proxy written in golang, including direct, GAE, SOCKS5, and SSH fetcher.


Status  [![Build Status](https://travis-ci.org/justmao945/mallory.svg?branch=master)](https://travis-ci.org/justmao945/mallory)
=============
* Support direct fetcher that spawn links from the running machine.
* Support GAE fetcher, only support http and https with port 443. In this mode we need to deploy fake certificates to forward https requests...
* Support SOCKS5 proxy fetcher, aka SOCKS5 to HTTP proxy translator.
* Support SSH fetcher, aka HTTP proxy via SSH tunnel.
* Simple PAC file server.

Installation
=============
```sh
go get github.com/justmao945/mallory/cmd/mallory
```

Engines
=============
### Direct
```sh
# This is the default mode, that spawns connections from the running machine.
# Now we have the HTTP proxy on port 1315
mallory
2014/04/12 01:56:33 Listen and serve HTTP proxy on 127.0.0.1:1315
2014/04/12 01:56:33 	Engine: direct
```

### GAE
```sh
# This engine spawns connections from the remote Google Application Engine.

# copy config to the default work dir
mkdir ~/.mallory && cp cfg/* ~/.mallory

# before start the proxy server, we'd better upload the GAE remote application
# for details see https://appengine.google.com
cd mallory/

# put your own app id into app.yaml
vim gae_server/app.yaml

# deploy it with go_appengine, https://developers.google.com/appengine/downloads
goapp deploy gae_server/

# this mode need to use the fake CA, default include crt and key
mallory -engine=gae -remote=https://your-app-id.appspot.com

# or generate the private key and sign the Root CA by yourself
openssl genrsa -out key 2048
openssl req -new -x509 -days 365 -key key -out crt
```

### SOCKS
```sh
# This engine spawns connections from the remote SOCKS proxy server,
# a simple way to translate SOCKS proxy to HTTP proxy.
# Assume we have a SOCKS5 proxy server at localhost and listen on port 1314
# Now we have the HTTP proxy on port 1315
mallory -engine=socks -remote=socks5://localhost:1314
2014/06/19 16:39:05 Starting...
2014/06/19 16:39:05 Listen and serve HTTP proxy on 127.0.0.1:1315
2014/06/19 16:39:05 	Engine: socks
2014/06/19 16:39:05 	Remote SOCKS proxy server: socks5://localhost:1314
```

### SSH
```sh
# This engine spawns connections from the remote SSH server, similar to the ssh -D command.
# The difference between them is that:
#   ssh -D  ==> SOCKS proxy
#   mallory ==> HTTP proxy
# Assume we have a ssh server on linode:22
# Now we have the HTTP proxy on port 1315
mallory -engine=ssh -remote=ssh://linode:22
2014/06/19 16:45:12 Starting...
2014/06/19 16:45:13 Listen and serve HTTP proxy on 127.0.0.1:1315
2014/06/19 16:45:13 	Engine: ssh
2014/06/19 16:45:13 	Remote SSH server: ssh://linode:22

# Add username, password and custom port 1234
mallory -engine=ssh -remote=ssh://user:password@linode:1234
```


TODO
=============
* Add test
* ....


References
=============
* [goproxy][1]
* [mitmproxy][2]
* [goagent][3]
 

[1]: https://github.com/elazarl/goproxy
[2]: http://mitmproxy.org/
[3]: https://github.com/goagent
