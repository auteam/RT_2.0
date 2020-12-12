package pages

import (
	"encoding/json"
	"net/http"
)

type Pages interface {
	Admin()
	MainPage()
}

// Admin ...
func Admin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// ctx := r.Context()
		// u := ctx.(Get)
		// w.Write([]byte(fmt.Sprintf("Welcome to admin room, %v!\n", u.Email)))
		// respond(w, r, http.StatusOK, nil)
	}
}

func respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
