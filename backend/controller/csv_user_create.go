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
func (s *server) UserFromCSVCreate() http.HandlerFunc {
	type request struct {
		Email    string
		Password string
		Name     string
	}
	type record struct {
		Accept  int
		Discard int
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		u := &model.User{}

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

			u = &model.User{
				Email:    req.Email,
				Password: req.Password,
				Name:     req.Name,
			}
			fmt.Println(u)
			if x, err := s.store.User().FindByEmail(u.Email); err == nil {
				fmt.Println("Dublicate")
				fmt.Println(x)
				records.Discard++
				i++
				continue
			}
			fmt.Println(u)
			// Create user
			if err := s.store.User().Create(u); err != nil {
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
