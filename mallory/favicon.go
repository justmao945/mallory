package mallory

import (
	"net/http"
)

// Provide a simple service for favicon
type ServiceFavicon struct{}

// main handler
func (*ServiceFavicon) Serve(s *Session) {
	s.Info("RESPONSE %s", StatusText(http.StatusOK))
}

// return "/favicon.ico"
func (*ServiceFavicon) Path() string {
	return "/favicon.ico"
}
