package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//CreateTopology ...
func (s *server) CreateTopology() http.HandlerFunc {
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
		top := fmt.Sprintf("{\"Name\":\"%v\",\"Devices\":\"\",\"Lines\":\"\",\"Keys\":\"\"}", name)
		_, err = s.store.User().CreateTopology(top, name)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		w.Write([]byte(fmt.Sprintf("{\"status\":\"OK\"}")))
		s.respond(w, r, http.StatusOK, nil)
	}
}
