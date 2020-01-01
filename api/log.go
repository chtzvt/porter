package api

import (
	"fmt"
	"net/http"
	"time"
)

func (s *Server) log(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%v API %v %v %v %v %v \"%v\"\n", time.Now(), r.RemoteAddr, r.Host, r.Proto, r.Method, r.RequestURI, r.UserAgent())
		h(w, r)
	}
}
