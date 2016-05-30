package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
	"net"
	"net/http"
	"strings"
)

type IPAuth struct {
	*AbstractAuth
}

func (a IPAuth) SuccessLogLevel() string {
	return a.GetLogLevel("auth.ip_auth.success_log_level", "Info")
}

func (a IPAuth) FailureLogLevel() string {
	return a.GetLogLevel("auth.ip_auth.failure_log_level", "Warn")
}

func (a IPAuth) Authenticate(r *http.Request, ps httprouter.Params) bool {
	line := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(line, ",")[0])
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	permitIps := viper.GetStringSlice("auth.ip_auth.ips")
	for _, v := range permitIps {
		if v == ip {
			return true
		}
	}
	return false
}

func (a IPAuth) HandleAuthError(w http.ResponseWriter, r *http.Request, ps httprouter.Params) bool {
	http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	return true
}
