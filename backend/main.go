package main

import (
	"flag"
	"log"

	//"net/http"
	"./controller"
	//"github.com/gorilla/mux"
	//"github.com/BurntSushi/toml"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")

}

func main() {
	flag.Parse()

	if err := controller.Start(); err != nil {
		log.Fatal(err)
	}
	// r := mux.NewRouter()
	// r.HandleFunc("/register", controller.RegisterHandler()).Methods("POST")
	// //r.HandleFunc("/login", LoginHandler).Methods("POST")

	// log.Fatal(http.ListenAndServe(Config.BindAddr, r))
}
