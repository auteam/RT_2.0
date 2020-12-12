package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//CreateModule ...
func (s *server) GetTopology() http.HandlerFunc {
	type request struct {
		Name string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		_, top, err := s.store.User().GetTopology(req.Name)
		if err != nil {
			s.error(w, r, http.StatusOK, err)
			return
		}

		w.Write([]byte(fmt.Sprintf("%v", top)))
		s.respond(w, r, http.StatusOK, nil)
	}
}
