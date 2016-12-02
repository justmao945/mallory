package mallory

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"sync"
	"syscall"

	"gopkg.in/fsnotify.v1"
)

// Memory representation for mallory.json
type ConfigFile struct {
	// private file file
	PrivateKey string `json:"id_rsa"`
	// local addr to listen and serve, default is 127.0.0.1:1315
	LocalSmartServer string `json:"local_smart"`
	// local addr to listen and serve, default is 127.0.0.1:1316
	LocalNormalServer string `json:"local_normal"`
	// remote addr to connect, e.g. ssh://user@linode.my:22
	RemoteServer string `json:"remote"`
	// direct to proxy dial timeout
	ShouldProxyTimeoutMS int `json:"should_proxy_timeout_ms"`
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
func (self *ConfigFile) Blocked(host string) bool {
	i := sort.SearchStrings(self.BlockedList, host)
	return i < len(self.BlockedList) && self.BlockedList[i] == host
}

// Provide global config for mallory
type Config struct {
	// file path
	Path string
	// config file content
	File *ConfigFile
	// File wather
	Watcher *fsnotify.Watcher
	// mutex for config file
	mutex  sync.RWMutex
	loaded bool
}

func NewConfig(path string) (self *Config, err error) {
	// watch config file changes
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}

	self = &Config{
		Path:    os.ExpandEnv(path),
		Watcher: watcher,
	}
	err = self.Load()
	return
}

func (self *Config) Reload() (err error) {
	file, err := NewConfigFile(self.Path)
	if err != nil {
		L.Printf("Reload %s failed: %s\n", self.Path, err)
	} else {
		L.Printf("Reload %s\n", self.Path)
		self.mutex.Lock()
		self.File = file
		self.mutex.Unlock()
	}
	return
}

// reload config file
func (self *Config) Load() (err error) {
	if self.loaded {
		panic("can not be reload manually")
	}
	self.loaded = true

	// first time to load
	L.Printf("Loading: %s\n", self.Path)
	self.File, err = NewConfigFile(self.Path)
	if err != nil {
		return
	}

	// Watching the whole directory instead of the individual path.
	// Because many editors won't write to file directly, they copy
	// the original one and rename it.
	err = self.Watcher.Add(filepath.Dir(self.Path))
	if err != nil {
		return
	}

	go func() {
		for {
			select {
			case event := <-self.Watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write && event.Name == self.Path {
					self.Reload()
				}
			case err := <-self.Watcher.Errors:
				L.Printf("Watching failed: %s\n", err)
			}
		}
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP)
	go func() {
		for s := range sc {
			if s == syscall.SIGHUP {
				self.Reload()
			}
		}
	}()

	return
}

// test whether host is in blocked list or not
func (self *Config) Blocked(host string) bool {
	self.mutex.RLock()
	blocked := self.File.Blocked(host)
	self.mutex.RUnlock()
	return blocked
}
