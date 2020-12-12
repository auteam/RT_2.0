package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"strconv"
)

//AllStand ...
func (s *server) SetTime() http.HandlerFunc {
	type request struct {
		Name string
		TimeEnd string
		TimeZone string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		datetime, err := time.Parse(time.RFC3339, req.TimeEnd)
		if err != nil {
			s.error(w, r, http.StatusTeapot, err)
			return
		}

		location, err := strconv.Atoi(req.TimeZone)
		if err != nil {
			s.error(w, r, http.StatusTeapot, err)
			return
		}

		datetime = datetime.Add(-(time.Duration(location) * time.Hour))

		timeRFC := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d.00Z",
			datetime.Year(), datetime.Month(), datetime.Day(),
			datetime.Hour(), datetime.Minute(), datetime.Second())


		_, jsonTopology, err := s.store.User().GetTopology(req.Name)
		if err != nil {
			s.error(w, r, http.StatusTeapot, err)
			return
		}

		err = s.store.User().SaveTopology(fmt.Sprintf("%s%s\"}", jsonTopology[0:len(jsonTopology) - 25], timeRFC), req.Name)
		if err != nil {
			s.error(w, r, http.StatusTeapot, err)
			return
		}

		w.Write([]byte("{\"status\":\"OK\"}"))
		s.respond(w, r, http.StatusOK, nil)
	}
}
