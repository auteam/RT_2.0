package model

// Stand ...
type Stand struct {
	ID         int    `json:"id"`
	Email      string `json:"Email"`
	Address    string `json:"Address"`
	Digipass   string
	Datacenter string
	Digiuser   string
	Esxipass   string
	Esxiuser   string
	Module     string
	Digi       string
	PortT      string
	Port       map[string]interface{}
}

// ValidateStand ...
func (u *Stand) ValidateStand() error {

	return nil
}
