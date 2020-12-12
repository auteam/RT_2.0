package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"../vmrc"
)

//Topology
func (s *server) Topology() http.HandlerFunc {
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

		name := []string{}
		link := []string{}
		ssh := []string{}
		for i := range stand {
			fmt.Println("----- %s", stand[i])
			if stand[i].Address != "" && stand[i].Datacenter != "" && stand[i].Esxipass != "" && stand[i].Esxipass != "" {
				name, link, err = vmrc.VMRC(stand[i].Datacenter, stand[i].Address, stand[i].Esxiuser, stand[i].Esxipass)
				if err != nil {
					fmt.Println(err)
				}
			// TEMP CODE FROM 05.12.2020 \/\/\/
			} else if stand[i].Address != "" && stand[i].Datacenter == "" && stand[i].Esxipass == "" && stand[i].Esxipass == "" {
				name = []string{"SRV1", "PC1", "PC2"}
				link = []string{fmt.Sprintf("vnc://%s:32769", stand[i].Address), fmt.Sprintf("vnc://%s:32772", stand[i].Address), fmt.Sprintf("vnc://%s:32779", stand[i].Address)}
			}
			// TEMP CODE FROM 05.12.2020 /\/\/\
			fmt.Println(stand[i].PortT)
			ssh = strings.Split(stand[i].PortT, ",")
		}

		str := "{"
		for i := 0; i < len(link); i++ {
			str = str + "\"" + name[i] + "\":\"" + link[i] + "\",\n"
		}
		str = str + "}"

		for i := range strings.Split(res2, " ") {
			temp := strings.Split(res2, " ")[i]
			temp2 := device[temp].(map[string]interface{})
			k := fmt.Sprintf("%v", temp2["vm"])
			for j := 0; j < len(link); j++ {
				if name[j] == k {
					temp2["link"] = link[j]
				}
			}
			for j := range ssh {
				if strings.Split(ssh[j], ":")[0] == temp2["vm"] {
					port := strings.Split(ssh[j], ":")
					telnet_or_ssh := strings.Split(port[1], "/")
					if telnet_or_ssh[0] == "ssh" {
						temp2["link"] = fmt.Sprintf("ssh://%v:%s", stand[0].Digi, telnet_or_ssh[1])
					} else if telnet_or_ssh[0] == "telnet" {
						temp2["link"] = fmt.Sprintf("telnet://%v:%s", stand[0].Digi, telnet_or_ssh[1])
					}
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
