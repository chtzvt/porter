package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *Server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, s.Config.HTTP.IndexFile)
	}
}

func (s *Server) handleListDoors() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setResponseHeaders(w)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(s.DoorMap)
	}
}

func (s *Server) handleDoorState() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d, err := s.getDoorFromPath(r)
		if err != nil {
			http.Error(w, NewStatusString("error", err.Error()), http.StatusBadRequest)
			return
		}

		setResponseHeaders(w)
		w.WriteHeader(http.StatusOK)

		_ = json.NewEncoder(w).Encode(d)
	}
}

func (s *Server) handleDoorOpen() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d, err := s.getDoorFromPath(r)
		if err != nil {
			http.Error(w, NewStatusString("error", err.Error()), http.StatusBadRequest)
			return
		}

		setResponseHeaders(w)

		stat := NewStatus()

		if err = d.Open(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			stat.Set("error", err.Error())
		} else {
			w.WriteHeader(http.StatusOK)
			stat.Set("OK", fmt.Sprintf("%s opening", d.Name))
		}

		_ = json.NewEncoder(w).Encode(stat)
	}
}

func (s *Server) handleDoorClose() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d, err := s.getDoorFromPath(r)
		if err != nil {
			http.Error(w, NewStatusString("error", err.Error()), http.StatusBadRequest)
			return
		}

		setResponseHeaders(w)

		stat := NewStatus()

		if err = d.Close(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			stat.Set("error", err.Error())
		} else {
			w.WriteHeader(http.StatusOK)
			stat.Set("OK", fmt.Sprintf("%s closing", d.Name))
		}

		_ = json.NewEncoder(w).Encode(stat)
	}
}

func (s *Server) handleDoorTrip() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d, err := s.getDoorFromPath(r)
		if err != nil {
			http.Error(w, NewStatusString("error", err.Error()), http.StatusBadRequest)
			return
		}

		setResponseHeaders(w)

		stat := NewStatus()

		if err = d.Trip(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			stat.Set("error", err.Error())
		} else {
			w.WriteHeader(http.StatusOK)
			stat.Set("OK", fmt.Sprintf("%s tripping", d.Name))
		}

		_ = json.NewEncoder(w).Encode(stat)
	}
}

func (s *Server) handleDoorLock() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d, err := s.getDoorFromPath(r)
		if err != nil {
			http.Error(w, NewStatusString("error", err.Error()), http.StatusBadRequest)
			return
		}

		setResponseHeaders(w)

		stat := NewStatusString("OK", fmt.Sprintf("%s locked", d.Name))

		d.Lock()

		_ = json.NewEncoder(w).Encode(stat)
	}
}

func (s *Server) handleDoorUnlock() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d, err := s.getDoorFromPath(r)
		if err != nil {
			http.Error(w, NewStatusString("error", err.Error()), http.StatusBadRequest)
			return
		}

		setResponseHeaders(w)

		stat := NewStatusString("OK", fmt.Sprintf("%s unlocked", d.Name))

		d.Unlock()

		_ = json.NewEncoder(w).Encode(stat)
	}
}

func setResponseHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Server", "porter")
}
