package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

//Topology
func (s *server) VNCTopology() http.HandlerFunc {
	type request struct {
		Champ  string `json:"champ"`
		Module string `json:"module"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		champ, err := s.store.User().GetChamp(req.Champ, s.email)

		if champ.Issue {
			if err := s.store.User().IssueStand(req.Champ, s.email); err != nil {
				w.Write([]byte(fmt.Sprintf("%v\n", err)))
				return
			}
		} else {
			//w.Write([]byte(fmt.Sprintf("%v\n", "Admin don't gave stand")))
			//return
		}
		stand, err := s.store.User().GetStand(req.Champ, s.email, req.Module)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("%v\n", err)))
			return
		}
		data, err := s.store.User().GetTopologyForUser(strings.ToLower(req.Champ), strings.ToUpper(req.Module))
		if err != nil {
			w.Write([]byte(fmt.Sprintf("%v\n", err)))
			return
		}
		var jsonStruct map[string]interface{}
		json.Unmarshal([]byte(``+string(data)+``), &jsonStruct)
		ttt := strings.Split(champ.Moduls[0], ";")
		for i := range ttt {

			start := time.Now()
			t, err := time.Parse(time.RFC3339, ttt[i])
			if err != nil {
			}

			tt := start.Format(time.RFC3339)
			tstart, err := time.Parse(time.RFC3339, tt)
			if err != nil {
			}

			fmt.Println(tstart, t, ttt[i], tt)

			g1 := start.Before(t)

			if g1 && (i == 1 || i == 3) {
				jsonStruct["Time"] = ttt[i]
				break
			} else {
				fmt.Println("err")
			}

		}
		keys := fmt.Sprintf("%v", jsonStruct["Keys"])
		res1 := strings.Replace(keys, "]", "", 1)
		res2 := strings.Replace(res1, "[", "", 1)
		device := jsonStruct["Devices"].(map[string]interface{})

		ssh := []string{}
		for i := range stand {
			fmt.Println(stand[i].PortT)
			ssh = strings.Split(stand[i].PortT, ",")
		}

		for i := range strings.Split(res2, " ") {
			temp := strings.Split(res2, " ")[i]
			temp2 := device[temp].(map[string]interface{})
			for j := range ssh {
				n := strings.Split(ssh[j], ":")[0]
				if strings.Split(n, "+")[0] == "T" && n == temp2["vm"] {
					temp2["link"] = fmt.Sprintf("telnet://%s:%s", stand[0].Digi, strings.Split(ssh[j], ":")[1])
				}
				if strings.Split(n, "+")[0] == "V" && n == temp2["vm"] {
					temp2["link"] = fmt.Sprintf("vnc://%s:%s", stand[0].Digi, strings.Split(ssh[j], ":")[1])
				}
			}
		}

		err = json.NewEncoder(w).Encode(jsonStruct)

		// if err != nil {
		// 	s.error(w, r, http.StatusBadRequest, err)
		// 	return
		// }
		w.WriteHeader(http.StatusOK)
		//s.respond(w, r, http.StatusOK, nil)
	}
}
