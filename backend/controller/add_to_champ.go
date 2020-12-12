package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"../model"
)

//UserFromCSV ...
func (s *server) AddToChamp() http.HandlerFunc {
	type request struct {
		Email  string
		Champ  string
		Module string
		Time   []string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		//mod := ["2020-11-05T21:25:20+05:00;2020-11-05T21:00:20+05:00;2020-11-05T21:25:20+05:00;2020-11-06T23:25:20+05:00%2020-11-05T21:25:20+05:00;2020-11-05T21:00:20+05:00;2020-11-05T21:25:20+05:00;2020-11-06T23:25:20+05:00%2020-11-05T21:25:20+05:00;2020-11-05T21:00:20+05:00;2020-11-05T21:25:20+05:00;2020-11-06T23:25:20+05:00"]
		var mod []string
		mod = append(mod, "2020-11-05T21:25:20+05:00;2020-11-05T21:00:20+05:00;2020-11-05T21:25:20+05:00;2020-11-06T23:25:20+05:00")
		mod = append(mod, "2020-11-05T21:25:20+05:00;2020-11-05T21:00:20+05:00;2020-11-05T21:25:20+05:00;2020-11-06T23:25:20+05:00")
		mod = append(mod, "2020-11-05T21:25:20+05:00;2020-11-05T21:00:20+05:00;2020-11-05T21:25:20+05:00;2020-11-06T23:25:20+05:00")

		u := &model.Champs{
			Email:  req.Email,
			Module: req.Module,
			Moduls: mod,
		}
		fmt.Println(u, req.Champ)
		if err := s.store.User().AddToChamp(u, req.Champ); err != nil {
			s.error(w, r, http.StatusTeapot, err)
			return
		}

		w.Write([]byte(fmt.Sprintf("{\"status\":\"OK\"}")))
		s.respond(w, r, http.StatusOK, nil)
	}
}
