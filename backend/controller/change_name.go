package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

//UserFromCSV ...
func (s *server) ChangeName() http.HandlerFunc {
	type request struct {
		Email string
		Name  string
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

			if err := s.store.User().ChangeName(req.Email, req.Name); err != nil {
				s.error(w, r, http.StatusTeapot, err)
				continue
			}
			i++
		}
		w.Write([]byte(fmt.Sprintf("{\"status\":\"OK\"}")))
		s.respond(w, r, http.StatusOK, nil)
	}
}
