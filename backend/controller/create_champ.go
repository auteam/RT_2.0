package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

//CreateTopology ...
func (s *server) CreateChamp() http.HandlerFunc {
	type request struct {
		Champ string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if req.Champ == "" {
			s.error(w, r, http.StatusBadRequest, errors.New("Wrong championship name"))
			return
		}
		err := s.store.User().CreateChamp(strings.ToLower(req.Champ))
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil)
	}
}
