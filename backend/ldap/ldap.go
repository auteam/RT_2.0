package ldap

import (
	"fmt"
	"strings"

	auth "github.com/korylprince/go-ad-auth"
)

// UserRepository ...
type ldap struct {
}

const (
	filterDN = "(&(objectClass=person)(CN=Laberskaya,DC=au,DC=team)(|(sAMAccountName={username})(mail={username})))"
	baseDN   = "OU=Laberskaya,DC=example,DC=com"

	loginUsername = "tboerger"
	loginPassword = "password"
)

// Create ...
func Create() bool {

	Config := &auth.Config{
		Server:   "ad.au.team",
		Port:     389,
		BaseDN:   "OU=OVPN,DC=au,DC=team",
		Security: auth.SecurityInsecureStartTLS,
	}

	username := "vzdornov.ry"
	password := "P@ssw0rd"

	status, err := auth.Authenticate(Config, username, password)
	if err != nil {
		//handle err
		return status
	}

	if !status {
		//handle failed authentication
		return status
	}
	upn, err := Config.UPN(username)
	Conn, err := Config.Connect()
	st, err := Conn.Bind(upn, password)
	if err != nil {
		fmt.Println(st, err)
	}
	search, err := Conn.Search("(cn=*)", []string{"cn", "name"}, 1000)
	if err != nil {
		fmt.Println(st, err)
	}
	for i := range search {
		fmt.Println(search[i].GetAttributeValues("cn"), search[i].GetAttributeValues("name"))
	}
	return status
}

func filter(needle string) string {
	res := strings.Replace(
		filterDN,
		"{username}",
		needle,
		-1,
	)

	return res
}
