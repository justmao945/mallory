package mallory

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sort"
)

// Memory representation for mallory.json
type ConfigFile struct {
	// private file file
	PrivateKey string `json:"id_rsa"`
	// local addr to listen and serve, default is 127.0.0.1:1315
	LocalServer string `json:"local"`
	// remote addr to connect, e.g. ssh://user@linode.my:22
	RemoteServer string `json:"remote"`
	// blocked host list
	BlockedList []string `json:"blocked"`
}

// Load file from path
func NewConfigFile(path string) (self *ConfigFile, err error) {
	self = &ConfigFile{}
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	err = json.Unmarshal(buf, self)
	if err != nil {
		return
	}
	self.PrivateKey = os.ExpandEnv(self.PrivateKey)
	sort.Strings(self.BlockedList)
	return
}

// test whether host is in blocked list or not
func (self *ConfigFile) Contain(host string) bool {
	i := sort.SearchStrings(self.BlockedList, host)
	return i < len(self.BlockedList) && self.BlockedList[i] == host
}

// Provide global config for mallory
type Config struct {
	// file path
	Path string
	// config file content
	File *ConfigFile
}

func NewConfig(path string) (self *Config, err error) {
	self = &Config{Path: os.ExpandEnv(path)}
	err = self.Load()
	return
}

// reload config file
func (self *Config) Load() (err error) {
	self.File, err = NewConfigFile(self.Path)
	return
}

// test whether host is in blocked list or not
func (self *Config) Contain(host string) bool {
	return self.File.Contain(host)
}
