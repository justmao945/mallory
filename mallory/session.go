package mallory

import (
	"log"
	"net/http"
	"path"
	"runtime"
	"strconv"
)

// A session is a proxy request
type Session struct {
	// Global server
	Server *Server
	// the unique ID start from 1
	ID int64
	// Copy from the http handler
	ResponseWriter http.ResponseWriter
	Request        *http.Request
}

func NewSession(s *Server, w http.ResponseWriter, r *http.Request) *Session {
	id, err := strconv.ParseInt(r.Header.Get("Mallory-Session"), 0, 64)
	if err != nil {
		id = s.NewID()
	}
	return &Session{
		Server:         s,
		ID:             id,
		ResponseWriter: w,
		Request:        r,
	}
}

func (self *Session) printf(format string, args ...interface{}) {
	br := []interface{}{self.ID, self.Server.CountAlive}
	log.Printf("[%03d/%02d] "+format+"\n", append(br, args...)...)
}

func (self *Session) printatf(ty, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	loc := []interface{}{path.Base(file), line}
	self.printf(ty+"%s:%d:"+format, append(loc, args...)...)
}

func (self *Session) Info(format string, args ...interface{}) {
	self.printf("INFO: "+format, args...)
}

func (self *Session) Warn(format string, args ...interface{}) {
	if self.Server.Env.Istty {
		self.printatf(CO_YELLOW+"WARN: ", format+CO_RESET, args...)
	} else {
		self.printatf("WARN: ", format, args...)
	}
}

func (self *Session) Error(format string, args ...interface{}) {
	if self.Server.Env.Istty {
		self.printatf(CO_RED+"ERRO: ", format+CO_RESET, args...)
	} else {
		self.printatf("ERRO: ", format, args...)
	}
}
