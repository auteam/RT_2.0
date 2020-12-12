package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"../model"
)

//UserFromCSV ...
func (s *server) AddToChampCSVCreate() http.HandlerFunc {
	type request struct {
		Email  string
		Module string
		Time   []string
	}
	type record struct {
		Accept  int
		Discard int
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		u := &model.Champs{}

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
		fmt.Println(champ)
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

			u = &model.Champs{
				Email:  req.Email,
				Module: req.Module,
				Moduls: req.Time,
			}
			fmt.Println(u, champ)
			err = s.store.User().FindByEmailChamp(u.Email, champ)
			if err == nil {
				fmt.Println("Duplicate")
				records.Discard++
				i++
				continue
			}
			if err := s.store.User().AddToChamp(u, champ); err != nil {
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
