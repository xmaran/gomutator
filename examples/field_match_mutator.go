//go:build exclude

package main

import (
	"fmt"

	"github.com/cmmaran/gomutator"
)

func main() {
	type credentials struct {
		Username, Password string
	}

	credStruct := credentials{
		Username: "admin",
		Password: "Master#123",
	}

	credMap := map[string]string{
		"username": "admin",
		"password": "admin",
	}

	m := gomutator.NewFieldMatchMutator()
	pm := &gomutator.PasswordDefaultMutator{}
	m.Hook().Add("Password", pm)
	m.Hook().Add("password", pm)

	m.Execute(&credStruct)
	fmt.Printf("%#v\n", credStruct)

	m.Execute(&credMap)
	fmt.Printf("%#v\n", credMap)
}
