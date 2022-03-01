package main

import (
	"net/url"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/contract"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nhttp"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nval"
	"strings"
)

func processServerConfig(config *contract.Config) *nhttp.ServerConfig {
	// Parse config
	secure := nval.ParseBooleanFallback(config.ServerSecure, false)
	trustProxy := nval.ParseBooleanFallback(config.ServerSecure, false)
	debug := nval.ParseBooleanFallback(config.ServerSecure, false)

	// Init server config
	s := nhttp.ServerConfig{
		ListenPort: config.Port,
		Secure:     secure,
		TrustProxy: trustProxy,
		Debug:      debug,
	}

	// Get base URL
	baseUrlStr, ok := nval.ParseString(config.ServerBaseUrl)
	if !ok || baseUrlStr == "" {
		log.Warn("server.base_url is not set")
		s.BaseUrl = nhttp.BuildUrl("localhost", s.ListenPort, s.BasePath)
	} else {
		tmp, err := url.Parse(baseUrlStr)
		if err != nil {
			log.Warn("failed to parse server.base_url = " + baseUrlStr)
			s.BaseUrl = nhttp.BuildUrl("localhost", s.ListenPort, s.BasePath)
		} else {
			s.BaseUrl = *tmp
			s.BasePath = tmp.Path
		}
	}

	// Normalize base path
	if !strings.HasSuffix(s.BasePath, "/") {
		s.BasePath += "/"
	}

	return &s
}
