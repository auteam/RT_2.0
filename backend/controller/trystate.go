package controller

import (
	"fmt"
	"net/http"
	"strings"
	"encoding/json"
	"strconv"
)

//UserFromCSV ...
func (s *server) ResetTryState() http.HandlerFunc {
	type request struct {
		ID     string
		Status *bool
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		dict := strings.Split(req.ID, "@")
		if len(dict) == 2 {
			err := s.store.User().UpdateTryStateByEmail(req.ID, req.Status)
			if err != nil {
				s.error(w, r, http.StatusTeapot, err)
				return
			}
		} else {
			dict := strings.Split(req.ID, "-")
			if len(dict) == 2 {
				dict1, _ := strconv.Atoi(dict[0])
				dict2, _ := strconv.Atoi(dict[1])
				for i := dict1; i <= dict2; i++ {
					intOne := strconv.Itoa(i)
					err := s.store.User().UpdateTryStateByID(intOne, req.Status)
					if err != nil {
						s.error(w, r, http.StatusTeapot, err)
						return
					}
				}
			} else {
				dict := strings.Split(req.ID, ",")
				if len(dict) >= 2 {
					for _, value := range dict {
						err := s.store.User().UpdateTryStateByID(value, req.Status)
						if err != nil {
							s.error(w, r, http.StatusTeapot, err)
							return
						}
					}
				}
			}
		}
		if req.Status == nil {
			w.Write([]byte(fmt.Sprintf("{\"id\": \"%s\",\"status\": null}", req.ID)))
		} else {
			w.Write([]byte(fmt.Sprintf("{\"id\": \"%s\",\"status\": %t}", req.ID, *req.Status)))
		}
	}
}
