package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"../vmrc"
)

//AllStand ...
func (s *server) Snapshot() http.HandlerFunc {
	type request struct {
		Champ  string
		Module string
		Device string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		stand, err := s.store.User().GetStand(req.Champ, s.email, req.Module)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("%v\n", err)))
			return
		}
		for i := range stand {
			err = vmrc.Snap(stand[i].Address, stand[i].Esxiuser, stand[i].Esxipass, stand[i].Datacenter, req.Device)
			if err != nil {
				w.Write([]byte(fmt.Sprintf("%v\n", err)))
				return
			}
		}
		w.Write([]byte(fmt.Sprintf("{\"status\":\"OK\"}")))
		s.respond(w, r, http.StatusOK, nil)
	}
}
