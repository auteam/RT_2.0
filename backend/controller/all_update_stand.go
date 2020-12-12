package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"../model"
)

//AllStand ...
func (s *server) AllUpdateStand() http.HandlerFunc {
	type request struct {
		Champ      string
		Datacenter string
		Address    string
		Digipass   string
		Digiuser   string
		Esxipass   string
		Esxiuser   string
		Module     string
		Digi       string
		Email      string
		Port       map[string]interface{}
		ID         int
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
			return
		}
		var jsonStruct map[string]interface{}
		json.Unmarshal([]byte(``+string(body)+``), &jsonStruct)
		i := 1
		for err == nil {
			jsonbody, err := json.Marshal(jsonStruct[strconv.Itoa(i)])
			if err != nil {
				fmt.Println(err)
				break
			}
			if err := json.Unmarshal(jsonbody, &req); req == nil {
				fmt.Println(err)
				break
			}

			//stand := fmt.Sprintf("%v", jsonStruct[strconv.Itoa(i)])
			//json.Unmarshal([]byte(``+stand+``), &req)

			st := &model.Stand{}
			st.Datacenter = req.Datacenter
			st.Address = req.Address
			st.Digipass = req.Digipass
			st.Digiuser = req.Digiuser
			st.Esxipass = req.Esxipass
			st.Esxiuser = req.Esxiuser
			st.Digi = req.Digi
			st.Email = req.Email
			temp := fmt.Sprintf("%v", req.Port)
			res1 := strings.Replace(temp, "]", "", 1)
			res2 := strings.Replace(res1, "map[", "", 1)
			res2 = strings.Replace(res2, " ", ",", 100)
			st.Port = req.Port
			st.PortT = res2
			st.ID = req.ID
			fmt.Println(req.Champ)
			err = s.store.User().UpdateStand(req.Champ, req.Module, st)
			if err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
			i++
		}
		w.Write([]byte(fmt.Sprintf("{\"status\":\"OK\"}")))
		w.WriteHeader(http.StatusOK)
		//s.respond(w, r, http.StatusOK, nil)
	}
}
