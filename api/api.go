package api

import (
	"errors"
	"net/http"
	"porter/config"
	"porter/door"
	"strings"
)

// Server represents an API Server, including its configuration parameters and lookup tables.
type Server struct {
	Config  *config.Config
	KeyMap  map[string][]string
	DoorMap map[string]*door.Door
}

// NewServer returns a new instance of an api.Server with the provided config.Config.
func NewServer(c *config.Config) *Server {
	server := new(Server)
	server.Config = c

	server.KeyMap = make(map[string][]string)
	server.genKeyMap()

	server.DoorMap = make(map[string]*door.Door)
	server.genDoorMap()

	return server
}

// Serve registers HTTP handlers for each path in the API, and starts the HTTP(S) server
func (s *Server) Serve() error {
	mux := http.NewServeMux()

	if len(s.Config.HTTP.IndexFile) != 0 {
		mux.Handle("/", s.handleIndex())
	} else {
		mux.Handle("/", http.RedirectHandler("https://github.com/ctrezevant/porter", http.StatusFound))
	}

	mux.Handle("/api/v1/list", s.checkAuth(s.handleListDoors()))

	mux.Handle("/api/v1/state/", s.checkAuth(s.handleDoorState()))

	mux.Handle("/api/v1/lock/", s.checkAuth(s.handleDoorLock()))
	mux.Handle("/api/v1/unlock/", s.checkAuth(s.handleDoorUnlock()))

	mux.Handle("/api/v1/open/", s.checkAuth(s.handleDoorOpen()))
	mux.Handle("/api/v1/close/", s.checkAuth(s.handleDoorClose()))
	mux.Handle("/api/v1/trip/", s.checkAuth(s.handleDoorTrip()))

	if s.Config.HTTP.TLSCertFile != "" && s.Config.HTTP.TLSKeyFile != "" {
		return http.ListenAndServeTLS(s.Config.HTTP.ListenAddr, s.Config.HTTP.TLSCertFile, s.Config.HTTP.TLSKeyFile, mux)
	} else {
		return http.ListenAndServe(s.Config.HTTP.ListenAddr, mux)
	}
}

// Generates a lookup table mapping door names to their index in s.Config.Doors
func (s *Server) genDoorMap() {
	for i, d := range s.Config.Doors {
		s.DoorMap[d.Name] = door.New(s.Config.Doors[i])
	}
}

// getDoorFromPath returns a door by name from an API request's path, or an error if one doesn't exist.
func (s *Server) getDoorFromPath(r *http.Request) (*door.Door, error) {
	name := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1 : len(r.URL.Path)]
	return s.getDoorByName(name)
}

// getDoorByName returns a door matching the provided name, or an error if one doesn't exist.
func (s *Server) getDoorByName(name string) (*door.Door, error) {
	if d, ok := s.DoorMap[name]; ok {
		return d, nil
	} else {
		return &door.Door{}, errors.New("door doesn't exist")
	}

}
