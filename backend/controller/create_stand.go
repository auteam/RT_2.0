package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//CreateModule ...
func (s *server) CreateStand() http.HandlerFunc {
	type request struct {
		Champ  string
		Module string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		id, err := s.store.User().CreateStand(req.Champ, req.Module)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		w.Write([]byte(fmt.Sprintf("{\"status\":\"OK\",\"id\":\"%v\"}", id)))
		s.respond(w, r, http.StatusOK, nil)
	}
}
