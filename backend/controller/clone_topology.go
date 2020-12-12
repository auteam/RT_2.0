package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//CreateModule ...
func (s *server) CloneTopology() http.HandlerFunc {
	type request struct {
		NameFrom string
		NameTo   string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		_, t, err := s.store.User().GetTopology(req.NameFrom)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		_, err = s.store.User().CreateTopology(t, req.NameTo)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		w.Write([]byte(fmt.Sprintf("{\"status\":\"OK\"}")))
		s.respond(w, r, http.StatusOK, nil)
	}
}
