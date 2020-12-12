package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//AllStand ...
func (s *server) RemoveStand() http.HandlerFunc {
	type request struct {
		Champ  string
		Module string
		ID     int
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		err := s.store.User().RemoveStand(req.Champ, req.ID)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		w.Write([]byte(fmt.Sprintf("{\"status\":\"OK\"}")))
		s.respond(w, r, http.StatusOK, nil)
	}
}
