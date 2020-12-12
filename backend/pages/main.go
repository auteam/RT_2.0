package pages

import (
	"fmt"
	"net/http"
)

// Main ...
func MainPage(e string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf("Welcome to main page, %v!", e)))
		respond(w, r, http.StatusOK, nil)
	}
}
