package mallory

import (
	"log"
	"net/http"
	"path"
	"runtime"
)

// A session is a proxy request
type Session struct {
	// Global config
	Env *Env
	// the unique ID start from 1
	ID int64
	// Copy from the http handler
	ResponseWriter http.ResponseWriter
	Request        *http.Request
}

func NewSession(e *Env, id int64, w http.ResponseWriter, r *http.Request) *Session {
	return &Session{
		Env:            e,
		ID:             id,
		ResponseWriter: w,
		Request:        r,
	}
}

func (self *Session) printf(format string, args ...interface{}) {
	log.Printf("[%03d] "+format+"\n", append([]interface{}{self.ID}, args...)...)
}

func (self *Session) printatf(ty, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	self.printf(ty+"%s:%d:"+format, append([]interface{}{path.Base(file), line}, args...)...)
}

func (self *Session) Info(format string, args ...interface{}) {
	self.printf("INFO: "+format, args...)
}

func (self *Session) Warn(format string, args ...interface{}) {
	if self.Env.Istty {
		self.printatf(CO_YELLOW+"WARN: ", format+CO_RESET, args...)
	} else {
		self.printatf("WARN: ", format, args...)
	}
}

func (self *Session) Error(format string, args ...interface{}) {
	if self.Env.Istty {
		self.printatf(CO_RED+"ERRO: ", format+CO_RESET, args...)
	} else {
		self.printatf("ERRO: ", format, args...)
	}
}
