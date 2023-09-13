[![GoDoc](https://godoc.org/github.com/cmmaran/gomutator?status.svg)](https://godoc.org/github.com/cmmaran/gomutator)
[![Build Status](https://travis-ci.org/cmmaran/gomutator.svg)](https://travis-ci.org/cmmaran/gomutator)
[![Go Report Card](https://goreportcard.com/badge/github.com/cmmaran/gomutator)](https://goreportcard.com/report/github.com/cmmaran/gomutator)

# gomutator

## Description

gomutator provides functionality to modify the value of struct fields or map keys based on the matched key using the contract method.

For each field in the struct or map, it will visit them one by one and check if their name matches the input hook key. If there is a match, it will call the Mutate method of the hook contract. It is important to note that this package can only be used on the addressable value of the struct and map. Additionally, it can only be used on fields that have been explicitly exposed outside of the package.

## Installation

```
go get github.com/cmmaran/gomutator
```

Dependencies :
* [github.com/stretchr/testify/assert](https://github.com/stretchr/testify#assert-package)

## Usage

### Example

```go

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

```

Output:
```
main.credentials{Username:"admin", Password:"********"}
map[string]string{"password":"********", "username":"admin"}
```