package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/TakatoshiMaeda/kinu/logger"
	"net/http"
	"os"
	"reflect"
	"strings"
)

var (
	authenticationStack map[string][]Auth = map[string][]Auth{}
	authTypeRegistry map[string]reflect.Type = map[string]reflect.Type{
		"NoAuth":    reflect.TypeOf(NoAuth{}),
		"BasicAuth": reflect.TypeOf(BasicAuth{}),
	}
)

type Auth interface {
	Authenticate(r *http.Request, ps httprouter.Params) bool
	HandleAuthError(w http.ResponseWriter, r *http.Request, ps httprouter.Params) bool
	SuccessLogLevel() string
	FailureLogLevel() string
}

type AbstractAuth struct {
}

func (a *AbstractAuth) GetLogLevel(key string, def string) string {
	level := viper.GetString(key)
	if level == "" {
		return def
	}
	return level
}

func (a *AbstractAuth) SuccessLogLevel() string {
	return "Info"
}

func (a *AbstractAuth) FailureLogLevel() string {
	return "Warn"
}

type NoAuth struct {
	*AbstractAuth
}

func (a NoAuth) Authenticate(r *http.Request, ps httprouter.Params) bool {
	return true
}

func (a NoAuth) HandleAuthError(w http.ResponseWriter, r *http.Request, ps httprouter.Params) bool {
	return false
}

func createAuthInstance(name string) Auth {
	elem := reflect.New(authTypeRegistry[name]).Elem()
	return elem.Interface().(Auth)
}

func InitAuthStack() {
	viper.SetEnvPrefix("kinu")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigName("kinu_auth")
	viper.AddConfigPath(".")
	viper.ReadInConfig()
	if os.Getenv("KINU_DEBUG") == "1" {
		viper.Debug()
	}
	for _, k := range []string{"images", "upload"} {
		viper.SetDefault("resource." + k, []string{})
		for _, v := range viper.GetStringSlice("resource." + k) {
			authenticationStack[k] = append(authenticationStack[k], createAuthInstance(v))
		}
	}
}

func Authentication(resource string, h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if authenticate(w, r, ps, authenticationStack[resource], 0) {
			h(w, r, ps)
		}
	}
}

func authenticate(w http.ResponseWriter, r *http.Request, ps httprouter.Params, authStack []Auth, index int) bool {
	if len(authStack) - 1 < index {
		return true
	}
	auth := authStack[index]
	if auth.Authenticate(r, ps) {
		// auth success
		logger.Log(logger.WithFields(logrus.Fields{
			"path":   r.URL.Path,
			"params": r.URL.Query(),
			"method": r.Method,
			"auth": reflect.TypeOf(auth),
		}), auth.SuccessLogLevel(), "success")
		index += 1
		return authenticate(w, r, ps, authStack, index)
	}
	// auth error
	logger.Log(logger.WithFields(logrus.Fields{
		"path":   r.URL.Path,
		"params": r.URL.Query(),
		"method": r.Method,
		"auth": reflect.TypeOf(auth),
	}), auth.FailureLogLevel(), "failure")
	auth.HandleAuthError(w, r, ps)
	return false
}

