package mallory

import (
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os/user"
	"sync"

	"golang.org/x/crypto/ssh"
)

//
type SSH struct {
	// global config file
	Cfg *Config
	// connect URL
	URL *url.URL
	// SSH client
	Client *ssh.Client
	// SSH client config
	CliCfg *ssh.ClientConfig
	// direct fetcher
	Direct *Direct
	// only re-dial once
	sf Group
	l  sync.RWMutex
}

// Create and initialize
func NewSSH(c *Config) (self *SSH, err error) {
	self = &SSH{
		Cfg:    c,
		CliCfg: &ssh.ClientConfig{},
	}
	// e.g.  ssh://user:passwd@192.168.1.1:1122
	self.URL, err = url.Parse(c.File.RemoteServer)
	if err != nil {
		return
	}

	if self.URL.User != nil {
		self.CliCfg.User = self.URL.User.Username()
	} else {
		u, err := user.Current()
		if err != nil {
			return self, err
		}
		// u.Name is the full name, should not be used
		self.CliCfg.User = u.Username
	}

	self.CliCfg.HostKeyCallback = func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}

	// 1) try RSA keyring first
	for {
		id_rsa := c.File.PrivateKey
		pem, err := ioutil.ReadFile(id_rsa)
		if err != nil {
			L.Printf("ReadFile %s failed:%s\n", c.File.PrivateKey, err)
			break
		}
		signer, err := ssh.ParsePrivateKey(pem)
		if err != nil {
			L.Printf("ParsePrivateKey %s failed:%s\n", c.File.PrivateKey, err)
			break
		}
		self.CliCfg.Auth = append(self.CliCfg.Auth, ssh.PublicKeys(signer))
		// stop !!
		break
	}
	// 2) try password
	for {
		if self.URL.User == nil {
			break
		}
		if pass, ok := self.URL.User.Password(); ok {
			self.CliCfg.Auth = append(self.CliCfg.Auth, ssh.Password(pass))
		}
		// stop here!!
		break
	}

	if len(self.CliCfg.Auth) == 0 {
		//TODO: keyboard intercative
		err = errors.New("Invalid auth method, please add password or generate ssh keys")
		return
	}

	// first time to dial to remote server, make sure it is available
	self.Client, err = ssh.Dial("tcp", self.URL.Host, self.CliCfg)
	if err != nil {
		return
	}

	dial := func(network, addr string) (c net.Conn, err error) {
		self.l.RLock()
		cli := self.Client
		self.l.RUnlock()

		c, err = cli.Dial(network, addr)
		if err == nil {
			return
		}

		L.Printf("dial %s failed: %s, reconnecting ssh server %s...\n", addr, err, self.URL.Host)

		clif, err := self.sf.Do(network+addr, func() (interface{}, error) {
			return ssh.Dial("tcp", self.URL.Host, self.CliCfg)
		})
		if err != nil {
			L.Printf("connect ssh server %s failed: %s\n", self.URL.Host, err)
			return
		}
		cli = clif.(*ssh.Client)

		self.l.Lock()
		self.Client = cli
		self.l.Unlock()

		return cli.Dial(network, addr)
	}

	self.Direct = &Direct{
		Tr: &http.Transport{Dial: dial},
	}
	return
}

func (self *SSH) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	self.Direct.ServeHTTP(w, r)
}

func (self *SSH) Connect(w http.ResponseWriter, r *http.Request) {
	self.Direct.Connect(w, r)
}
