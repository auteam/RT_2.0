package model

//Champs ...
type Champs struct {
	ID       int      `json:"id"`
	Email    string   `json:"-"`
	Moduls   []string `json:"Moduls"`
	Time     string
	Module   string `json:"Module"`
	Standnum string `json:"-"`
	Issue    bool   `json:"-"`
}
