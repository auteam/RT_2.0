package controller

import (
	"fmt"
	"net/http"
	"strconv"
)

//SaveTopology ...
func (s *server) AllUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := s.store.User().AllUser()
		if err != nil {
			s.error(w, r, http.StatusOK, err)
			return
		}
		str := "{\n"
		for i := range u {
			if i != 0 {
				str += ",\n"
			}
			str += "\"" + strconv.Itoa(i+1) + "\":{\n"
			str += "\"id\":\"" + strconv.Itoa(u[i].ID) + "\",\n"
			str += "\"Email\":\"" + u[i].Email + "\",\n"
			str += "\"FIO\":\"" + u[i].Name + "\",\n"
			str += "\"Role\":\"" + u[i].Role + "\",\n"
			str += "\"Group\":\"" + u[i].Group + "\",\n"
			if u[i].TryState == nil {
				str += "\"TryState\": null\n"
			} else {
				str += "\"TryState\": " + strconv.FormatBool(*u[i].TryState) + "\n"
			}
			str += "}"
		}
		str += "}"
		w.Write([]byte(fmt.Sprintf(str)))
		s.respond(w, r, http.StatusOK, nil)
	}
}
