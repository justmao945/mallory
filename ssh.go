package mallory

import (
	"errors"
	"github.com/justmao945/mallory/ssh"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/user"
)

//
type EngineSSH struct {
	Env *Env
	URL *url.URL
	Cli *ssh.Client
	Cfg *ssh.ClientConfig
	Dir *EngineDirect
}

// Create and initialize
func CreateEngineSSH(e *Env) (self *EngineSSH, err error) {
	self = &EngineSSH{
		Env: e,
		Cfg: &ssh.ClientConfig{},
	}
	// e.g.  ssh://user:passwd@192.168.1.1:1122
	self.URL, err = url.Parse(e.Remote)
	if err != nil {
		return
	}

	if self.URL.User != nil {
		self.Cfg.User = self.URL.User.Username()
	} else {
		u, err := user.Current()
		if err != nil {
			return self, err
		}
		// u.Name is the full name, should not be used
		self.Cfg.User = u.Username
	}

	// 1) try RSA keyring first
	for {
		id_rsa := os.ExpandEnv("$HOME/.ssh/id_rsa")
		pem, err := ioutil.ReadFile(id_rsa)
		if err != nil {
			break
		}
		signer, err := ssh.ParsePrivateKey(pem)
		if err != nil {
			break
		}
		self.Cfg.Auth = append(self.Cfg.Auth, ssh.PublicKeys(signer))
		// stop !!
		break
	}
	// 2) try password
	for {
		if self.URL.User == nil {
			break
		}
		if pass, ok := self.URL.User.Password(); ok {
			self.Cfg.Auth = append(self.Cfg.Auth, ssh.Password(pass))
		}
		// stop here!!
		break
	}

	if len(self.Cfg.Auth) == 0 {
		//TODO: keyboard intercative
		err = errors.New("Invalid auth method, please add password or generate ssh keys")
		return
	}

	self.Cli, err = ssh.Dial("tcp", self.URL.Host, self.Cfg)
	if err != nil {
		return
	}

	dial := func(network, addr string) (c net.Conn, err error) {
		for i := 0; i < 3; i++ {
			c, err = self.Cli.Dial(network, addr)
			// We want to reconnect the network when disconnected.
			// FIXME: unexported net.errClosing
			if err != nil && err.Error() == "use of closed network connection" {
				continue
			}
			// do not reconnect when no error or other errors
			break
		}
		return c, err
	}

	self.Dir = &EngineDirect{
		Tr: &http.Transport{Dial: dial},
	}
	return
}

func (self *EngineSSH) Serve(s *Session) {
	self.Dir.Serve(s)
}

func (self *EngineSSH) Connect(s *Session) {
	self.Dir.Connect(s)
}
