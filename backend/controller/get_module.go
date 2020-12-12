package controller

import (
	"fmt"
	"net/http"
	"strings"
)

//CreateTopology ...
func (s *server) GetModule() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		module, err := s.store.User().GetModule()
		m := strings.Join(module[:], ",")
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if m == "" {
			w.Write([]byte(fmt.Sprintf("{\"status\":\"OK\",\"module\":{}}")))
		} else {
			w.Write([]byte(fmt.Sprintf("{\"status\":\"OK\",\"module\":\"%v\"}", m)))
		}
		s.respond(w, r, http.StatusOK, nil)
	}
}
