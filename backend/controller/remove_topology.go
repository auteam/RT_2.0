package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//RemoveTopology ...
func (s *server) RemoveTopology() http.HandlerFunc {
	type request struct {
		Name string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		err := s.store.User().RemoveTopology(req.Name)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		w.Write([]byte(fmt.Sprintf("{\"status\":\"OK\"}")))
		s.respond(w, r, http.StatusOK, nil)
	}
}
