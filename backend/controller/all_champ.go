package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

//CreateModule ...
func (s *server) AllChamp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ta, ch, err := s.store.User().AllChamp()
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		u, err := s.store.User().FindByEmail(s.email)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("{\"status\": \"400\",\"error\": \"Database access error\"}")))
			return
		}

		str := "{\n"
		var g []string

		for i := range ch {
			a := strings.Split(ch[i], ",")
			g = append(g, a[0])
		}

		if len(g) == 1 && len(g[0]) == 0 {
			str = str + "\"status\":\"CHAMPS_NOT_FOUND\",\n"
			str = str + "\"FIO\":\"" + u.Name + "\",\n"
			str = str + "\"Role\":\"" + u.Role + "\"\n"
			str = str + "}"
			w.Write([]byte(fmt.Sprintf("%v", str)))
			s.respond(w, r, http.StatusOK, nil)
			return
		}
		str = str + "\"Champs\":{\n"
		for i := range g {
			if i > 0 {
				str = str + ",\n"
			}
			str = str + "\"" + strconv.Itoa(i) + "\":{\n" +
				"\"name\":\"" + g[i] + "\",\n"

			module := strings.Split(ch[i], ",")[1:]
			str = str + "\"Moduls\":{\n"
			if len(module) == 1 && len(module[0]) == 0 {
				str = str + "}}"
				continue
			}
			t := strings.Split(ta[i], ",")
			for j := range module {
				if j > 0 {
					str = str + ",\n"
				}
				if t[j] == "" || err != nil {
					str = str + "\"" + module[j] + "\":{\"Topology\":false}"
				} else {
					id, err := strconv.Atoi(t[j])
					name, err := s.store.User().GetTopologyNameByID(id)
					if err != nil {
						continue
					}
					str = str + "\"" + module[j] + "\":{\"Topology\":\"" + name + "\"}"
				}
			}
			str = str + "}}"
		}
		if err == nil {
			str = str + "},\n"
		}
		str = str + "\"status\":\"OK\",\n"
		str = str + "\"FIO\":\"" + u.Name + "\",\n"
		str = str + "\"Role\":\"" + u.Role + "\"\n"
		str = str + "}"
		w.Write([]byte(fmt.Sprintf("%v", str)))

		s.respond(w, r, http.StatusOK, nil)
	}
}
