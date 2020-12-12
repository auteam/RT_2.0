package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"../vmrc"
)

//CreateModule ...
func (s *server) GetTicket() http.HandlerFunc {
	type request struct {
		Device string
		Champ  string
		Module string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var link string
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
			link, err = vmrc.Ticket(req.Device, stand[i].Datacenter, stand[i].Address, stand[i].Esxiuser, stand[i].Esxipass)
			if err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
		}
		//fmt.Println(link)
		w.Write([]byte(fmt.Sprintf("{\"link\":\"%s\"}", link)))
		s.respond(w, r, http.StatusOK, nil)
	}
}
