package controller

import (
	"fmt"
	"net/http"

	"../ldap"
)

//Админка
func (s *server) Admin() http.HandlerFunc {
	//return pages.Admin()
	return func(w http.ResponseWriter, r *http.Request) {
		status := ldap.Create()
		w.Write([]byte(fmt.Sprintf("Welcome %v!", status)))

		s.respond(w, r, http.StatusOK, nil)
	}
}
