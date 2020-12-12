package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

//AddToModule ...
func (s *server) AddToModule() http.HandlerFunc {
	type request struct {
		Email  string
		Champ  string
		Modules string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := s.store.User().AddToModule(strings.ToLower(req.Champ), req.Email, req.Modules); err != nil {
			s.error(w, r, http.StatusTeapot, err)
			return
		}

		w.Write([]byte(fmt.Sprintf("{\"status\":\"OK\"}")))
		s.respond(w, r, http.StatusOK, nil)
	}
}
