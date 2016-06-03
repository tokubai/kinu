package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"fmt"
)

func VersionHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, VERSION)
}
