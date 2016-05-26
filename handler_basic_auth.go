package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
	"net/http"
	"strings"
	"encoding/base64"
	"bytes"
)

type BasicAuth struct {
	*AbstractAuth
}

func (a BasicAuth) SuccessLogLevel() string {
	return a.GetLogLevel("auth.basic_auth.success_log_level", "Info")
}

func (a BasicAuth) FailureLogLevel() string {
	return a.GetLogLevel("auth.basic_auth.failure_log_level", "Warn")
}

func (a BasicAuth) Authenticate(r *http.Request, ps httprouter.Params) bool {
	const basicAuthPrefix string = "Basic "

	user := viper.GetString("auth.basic_auth.user")
	pass := viper.GetString("auth.basic_auth.pass")
	auth := r.Header.Get("Authorization")
	if strings.HasPrefix(auth, basicAuthPrefix) {
		payload, err := base64.StdEncoding.DecodeString(auth[len(basicAuthPrefix):])
		if err == nil {
			pair := bytes.SplitN(payload, []byte(":"), 2)
			if len(pair) == 2 &&
				bytes.Equal(pair[0], []byte(user)) &&
				bytes.Equal(pair[1], []byte(pass)) {
				return true
			}
		}
	}	
	return false
}

func (a BasicAuth) HandleAuthError(w http.ResponseWriter, r *http.Request, ps httprouter.Params) bool {
	w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	return true
}
