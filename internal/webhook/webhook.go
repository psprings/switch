package webhook

import (
	"net/http"
	"net/url"
)

// Config :
type Config struct {
	Header     http.Header
	Body       []byte
	Host       string
	RemoteAddr string
	Form       url.Values
	MessageID  string
}
