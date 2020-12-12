package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"../model"
)

//AllUpdateStand ...
func (s *server) UpdateStand() http.HandlerFunc {
	type request struct {
		Champ      string
		Datacenter string
		Address    string
		Digipass   string
		Digiuser   string
		Exsipass   string
		Exsiuser   string
		Module     string
		Digi       string
		Port       string
		ID         int
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		st := &model.Stand{}

		err := s.store.User().UpdateStand(req.Champ, req.Module, st)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		w.Write([]byte(fmt.Sprintf("{\"status\":\"OK\"}")))
		s.respond(w, r, http.StatusOK, nil)
	}
}
