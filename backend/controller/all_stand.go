package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

//AllStand ...
func (s *server) AllStand() http.HandlerFunc {
	type request struct {
		Champ  string
		Module string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		st, err := s.store.User().AllStand(req.Champ, req.Module)

		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		str := "{"
		for i := range st {
			if i > 0 {
				str = str + ","
			}
			if st[i].PortT != "" {
				st[i].PortT = "{\"" + strings.ReplaceAll(st[i].PortT, ":", "\":\"")
				st[i].PortT = strings.ReplaceAll(st[i].PortT, ",", "\",\"")
				st[i].PortT = st[i].PortT + "\"}"
			} else {
				st[i].PortT = "{}"
			}
			err := json.Unmarshal([]byte(st[i].PortT), &st[i].Port)
			if err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
			jsonData, err := json.Marshal(st[i])
			if err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
			str = str + "\"" + strconv.Itoa(i+1) + "\"" + ":" + string(jsonData)
		}
		str = str + "}"
		w.Write([]byte(fmt.Sprintf("%v", str)))
		s.respond(w, r, http.StatusOK, nil)
	}
}
