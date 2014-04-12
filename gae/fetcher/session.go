package fetcher

import (
	"appengine"
	"fmt"
	"net/http"
	"path"
	"runtime"
)

type Session struct {
	Ctx    appengine.Context
	Writer http.ResponseWriter
}

func NewSession(ctx appengine.Context, w http.ResponseWriter) *Session {
	return &Session{
		Ctx:    ctx,
		Writer: w,
	}
}

func (self *Session) HTTPError(format string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	s := fmt.Sprintf("%s:%d: "+format, append([]interface{}{path.Base(file), line}, args...)...)
	self.Ctx.Errorf(s + "\n")
	http.Error(self.Writer, s, http.StatusInternalServerError)
}
