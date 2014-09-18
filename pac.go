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
	Url *url.URL
}

// create and init
func CreateServicePAC(e *Env) (self *ServicePAC, err error) {
	self = &ServicePAC{Env: e}
	self.Url, err = url.Parse(e.PAC)
	return
}

// main handler, read file and response, please don't use cache
func (self *ServicePAC) Serve(s *Session) {
	pac, err = ioutil.ReadFile(self.Url.Path)
	if err != nil {
		s.ResponseWriter.WriteHeader(http.StatusNotFound)
		s.Error("RESPONSE %s", StatusText(http.StatusNotFound))
		return
	}

	s.ResponseWriter.Write(pac)
	s.Info("RESPONSE %s", StatusText(http.StatusOK))
}

// return "/pac"
func (self *ServicePAC) Path() string {
	return "/pac"
}
