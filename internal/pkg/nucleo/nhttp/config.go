package nhttp

import (
	"fmt"
	"net/url"
	"strings"
)

type ServerConfig struct {
	ListenPort int
	BasePath   string
	BaseUrl    url.URL
	Secure     bool
	TrustProxy bool
	Debug      bool
}

func (s *ServerConfig) GetListenPort() string {
	return fmt.Sprintf(":%d", s.ListenPort)
}

func (s *ServerConfig) GetHttpBaseUrl() string {
	u := s.BaseUrl
	if s.Secure {
		u.Scheme = "https"
	} else {
		u.Scheme = "http"
	}
	return u.String()
}

func (s *ServerConfig) GetWebSocketBaseUrl() string {
	u := s.BaseUrl
	if s.Secure {
		u.Scheme = "wss"
	} else {
		u.Scheme = "ws"
	}
	return u.String()
}

func (s *ServerConfig) GetBasePath() string {
	if s.BasePath == "" {
		return "/"
	}

	if s.BasePath == "/" || !strings.HasSuffix(s.BasePath, "/") {
		return s.BasePath
	}

	return strings.TrimRight(s.BasePath, "/")
}
