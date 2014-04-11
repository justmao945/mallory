package fetcher

import (
	"appengine"
	"fmt"
	"net/http"
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
	s := fmt.Sprintf(format, args...)
	self.Ctx.Errorf(s + "\n")
	http.Error(self.Writer, s, http.StatusInternalServerError)
}
