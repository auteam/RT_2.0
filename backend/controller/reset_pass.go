package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//SaveTopology ...
func (s *server) ResetPass() http.HandlerFunc {
	type request struct {
		Email    string
		Password string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		u, err := s.store.User().FindByEmail(req.Email)
		if err != nil {
			s.error(w, r, http.StatusOK, err)
			return
		}
		u.Password = req.Password
		err = s.store.User().ResetPass(u)
		if err != nil {
			s.error(w, r, http.StatusOK, err)
			return
		}

		w.Write([]byte(fmt.Sprintf("{\"status\": \"OK\"}")))
		s.respond(w, r, http.StatusOK, nil)
	}
}
