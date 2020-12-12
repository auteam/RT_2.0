package controller

import (
	"encoding/json"
	"net/http"
	"strings"
)

//CreateModule ...
func (s *server) DeleteModule() http.HandlerFunc {
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
		err := s.store.User().DeleteModule(strings.ToLower(req.Champ), strings.ToUpper(req.Module))
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil)
	}
}
