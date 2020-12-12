package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/lib/pq"
)

//Main
func (s *server) Main() http.HandlerFunc {
	//return pages.MainPage(s.email)
	return func(w http.ResponseWriter, r *http.Request) {
		//w.Write([]byte(fmt.Sprintf("{\"status\": \"%v\"}", s.email)))
		str := "{\n"
		u, err := s.store.User().FindByEmail(s.email)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("{\"status\": \"400\",\"error\": \"Database access error\"}")))
			return
		}
		g := strings.Split(u.Group, ",")
		fmt.Println(len(g[0]))
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
		    champ, err := s.store.User().GetChamp(g[i], s.email)
			if err != nil {
				if _, ok := err.(*pq.Error); ok {
					w.Write([]byte(fmt.Sprintf("{\"error\":\"%v\"}\n", errNoChamp)))
					//return
					break
				}
				//w.Write([]byte(fmt.Sprintf("%v\n", err)))
				if i > 0 {
                    str = str + ",\n"
                }
                str = str + "\"" + strconv.Itoa(i) + "\":{\n" +
                    "\"name\":\"" + g[i] + "\"\n}"
				//return
				break
			}

            if i > 0 {
				str = str + ",\n"
			}
			str = str + "\"" + strconv.Itoa(i) + "\":{\n" +
				"\"name\":\"" + g[i] + "\",\n"

			ml := strings.Split(champ.Module, ",")

			str = str + "\"Moduls\":{\n"
			// Cicle write modules
			for j := range ml {
				//m := strings.Split(champ.Moduls[j], ";")
				if j > 0 {
					str = str + ",\n"
				}
				//fmt.Println(m)
				//Format RFC3339 ... Example: 2020-09-30T12:30:00.00+05:00
				str = str + "\"" + ml[j] + "\":{\n" + // Name of Module
					"\"TimeStart\":\"" + "m[0]" + "\"," + // Time Start module
					"\"TimeEnd\":\"" + "m[1]" + "\"," + // Time End module
					"\"TimeStartPause\":\"" + "m[2]" + "\"," + // Time Start Pause module
					"\"TimeEndPause\":\"" + " m[3]" + "\"}\n" // Time End Pause module
			}
			str = str + "}}"
		}

		// 	str = str + "\"Moduls\":{\n"
		// 	// Cicle write modules
		// 	for j := range champ.Moduls {
		// 		m := strings.Split(champ.Moduls[j], ";")
		// 		if j > 0 {
		// 			str = str + ",\n"
		// 		}
		// 		fmt.Println(m)
		// 		//Format RFC3339 ... Example: 2020-09-30T12:30:00.00+05:00
		// 		str = str + "\"" + ml[j] + "\":{\n" + // Name of Module
		// 			"\"TimeStart\":\"" + m[0] + "\"," + // Time Start module
		// 			"\"TimeEnd\":\"" + m[1] + "\"," + // Time End module
		// 			"\"TimeStartPause\":\"" + m[2] + "\"," + // Time Start Pause module
		// 			"\"TimeEndPause\":\"" + m[3] + "\"}\n" // Time End Pause module
		// 	}
		// 	str = str + "}}"
		// }
		fmt.Println(err)
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
