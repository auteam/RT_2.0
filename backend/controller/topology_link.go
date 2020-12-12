package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//CreateTopology ...
func (s *server) TopologyLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
			return
		}
		var jsonStruct map[string]interface{}
		json.Unmarshal([]byte(``+string(body)+``), &jsonStruct)

		champ := fmt.Sprintf("%v", jsonStruct["Champ"])
		module := fmt.Sprintf("%v", jsonStruct["Module"])
		name := fmt.Sprintf("%v", jsonStruct["Name"])
		if name == "" {
			s.error(w, r, http.StatusBadRequest, errors.New("Wrong name"))
			return
		}
		id, _, err := s.store.User().GetTopology(name)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		fmt.Println(id)
		err = s.store.User().LinkTopology(champ, strings.ToUpper(module), id)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		w.Write([]byte(fmt.Sprintf("{\"status\":\"OK\"}")))
		s.respond(w, r, http.StatusOK, nil)
	}
}
