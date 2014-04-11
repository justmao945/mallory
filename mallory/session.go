package mallory

import (
	"log"
	"net/http"
)

// A session is a proxy request
type Session struct {
	// the unique ID start from 1
	ID int64
	// Copy from the http handler
	ResponseWriter http.ResponseWriter
	Request        *http.Request
}

func NewSession(id int64, w http.ResponseWriter, r *http.Request) *Session {
	return &Session{
		ID:             id,
		ResponseWriter: w,
		Request:        r,
	}
}

func (self *Session) printf(format string, args ...interface{}) {
	log.Printf("[%03d] "+format+"\n", append([]interface{}{self.ID}, args...)...)
}

func (self *Session) Info(format string, args ...interface{}) {
	self.printf("INFO: "+format, args...)
}

func (self *Session) Error(format string, args ...interface{}) {
	self.printf("ERRO: "+format, args...)
}
