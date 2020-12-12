package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//UpdateTopology ...
func (s *server) UpdateTopology() http.HandlerFunc {
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
		fmt.Println(champ, module)
		err = s.store.User().UpdateTopology(string(body), strings.ToLower(champ), strings.ToUpper(module))
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		w.Write([]byte(fmt.Sprintf("{\"status\":\"OK\"}")))
		//w.Write([]byte(fmt.Sprintf("Welcome %v!", s.email)))
		s.respond(w, r, http.StatusOK, nil)
	}
}
