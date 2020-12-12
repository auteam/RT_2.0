package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//SaveTopology ...
func (s *server) SaveTopology() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
			return
		}
		var jsonStruct map[string]interface{}
		json.Unmarshal([]byte(``+string(body)+``), &jsonStruct)

		name := fmt.Sprintf("%v", jsonStruct["Name"])
		if name == "" {
			s.error(w, r, http.StatusBadRequest, errors.New("Wrong name"))
			return
		}
		err = s.store.User().SaveTopology(string(body), name)
		if err != nil {
			s.error(w, r, http.StatusOK, err)
			return
		}

		w.Write([]byte(fmt.Sprintf("{\"status\": \"20\"}")))
		s.respond(w, r, http.StatusOK, nil)
	}
}
