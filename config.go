package mallory

import (
	"encoding/json"
	"io/ioutil"
	"sort"
)

type ConfigFile struct {
	// local addr to listen and serve, default is 127.0.0.1:1315
	Local string `json:"local"`
	// remote addr to connect, e.g. ssh://user@linode.my:22
	Remote string `json:"remote"`
	// blocked host list
	Blocked []string `json:"blocked"`
}

func NewConfigFile(path string) (err error, self *ConfigFile) {
	self = &ConfigFile{}
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	err = json.Unmarshal(buf, self)
	if err != nil {
		return
	}
	sort.Strings(self.Blocked)
}

func (self *ConfigFile) Search(host string) int {
	return self.Blocked.Search(host)
}

// Provide global config for mallory
type Config struct {
	// file path
	Path string
	// config
	File *ConfigFile
}

func NewConfig(path string) (error, *Config) {
	self := Config{Path: path}
	return self.Load(), self
}

func (self *Config) Load() (err error) {
	err, self.File = NewConfigFile(self.Path)
}
