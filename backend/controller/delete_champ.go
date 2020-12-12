package controller

import (
	"encoding/json"
	"net/http"
	"strings"
)

//CreateTopology ...
func (s *server) DeleteChamp() http.HandlerFunc {
	type request struct {
		Champ string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		err := s.store.User().DeleteChamp(strings.ToLower(req.Champ))
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil)
	}
}
