package mallory

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

// Provide a simple service for PAC file
type ServicePAC struct {
	// Global config
	Env *Env
	// the PAC file content
	PAC []byte
}

// create and init
func CreateServicePAC(e *Env) (self *ServicePAC, err error) {
	self = &ServicePAC{Env: e}

	url, err := url.Parse(e.PAC)
	if err != nil { // treat as a file path
		self.PAC, err = ioutil.ReadFile(e.PAC)
		return
	}

	if url.Scheme == "" || url.Scheme == "file" {
		self.PAC, err = ioutil.ReadFile(url.Path)
		return
	}

	resp, err := http.Get(e.PAC)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	_, err = resp.Body.Read(self.PAC)
	return
}

// main handler
func (self *ServicePAC) Serve(s *Session) {
	s.ResponseWriter.Write(self.PAC)
	s.Info("RESPONSE %s", StatusText(http.StatusOK))
}

// return "/pac"
func (self *ServicePAC) Path() string {
	return "/pac"
}
