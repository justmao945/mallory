package mallory

import (
	"github.com/justmao945/mallory/proxy"
	"net"
	"net/http"
	"net/url"
)

// translate SOCKS proxy to HTTP proxy
type EngineSOCKS struct {
	Env *Env
	Dir *EngineDirect
}

// Create and initialize
func CreateEngineSOCKS(e *Env) (self *EngineSOCKS, err error) {
	proxyURL, err := url.Parse(e.Remote)
	if err != nil {
		return
	}
	proxyDialer, err := proxy.FromURL(proxyURL, &http.DefaultTransport)
	if err != nil {
		return
	}

	dial := func(network, addr string) (net.Conn, error) {
		return proxyDialer.Dial(network, addr)
	}

	self = &EngineSOCKS{
		Env: e,
		Dir: &EngineDirect{Tr: &http.Transport{Dial: dial}},
	}
	return
}

func (self *EngineSOCKS) Serve(s *Session) {
	self.Dir.Serve(s)
}

func (self *EngineSOCKS) Connect(s *Session) {
	self.Dir.Connect(s)
}
