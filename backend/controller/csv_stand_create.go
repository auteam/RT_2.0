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

//UserFromCSV ...
func (s *server) StandFromCsvCreate() http.HandlerFunc {
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
	type record struct {
		Accept  int
		Discard int
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		records := &record{
			Accept:  0,
			Discard: 0,
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
			return
		}
		var jsonStruct map[string]interface{}
		json.Unmarshal([]byte(``+string(body)+``), &jsonStruct)

		champ := fmt.Sprintf("%v", jsonStruct["champ"])

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

			temp := fmt.Sprintf("%v", req.Port)
			res1 := strings.Replace(temp, "]", "", 1)
			res2 := strings.Replace(res1, "map[", "", 1)
			res2 = strings.Replace(res2, " ", ",", -1)

			st := &model.Stand{
				Datacenter: req.Datacenter,
				Address:    req.Address,
				Digipass:   req.Digipass,
				Digiuser:   req.Digiuser,
				Esxipass:   req.Esxipass,
				Esxiuser:   req.Esxiuser,
				Digi:       req.Digi,
				Email:      req.Email,
				Port:       req.Port,
				PortT:      res2,
				ID:         req.ID,
				Module:     req.Module,
			}

			fmt.Println(champ)
			fmt.Println(st)
			if err := s.store.User().CreateStandCSV(champ, st); err != nil {
				fmt.Println(err)
				records.Discard++
				i++
				continue
			}
			records.Accept++
			i++
		}

		w.Write([]byte(fmt.Sprintf("{\"status\":\"OK\",\"accept\":\"%v\",\"discard\":\"%v\"}", records.Accept, records.Discard)))
		s.respond(w, r, http.StatusOK, nil)
	}
}
