package gomutator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMutator(t *testing.T) {

	type CirclularCredentials struct {
		Head     *CirclularCredentials
		Username string
		Password string
	}

	cc := CirclularCredentials{}
	// test shouldn't crash as the mutator should be loop safe
	cc.Head = &cc
	cc.Username = "admin"
	cc.Password = "Master#123"

	cm := map[string]string{
		"username": "admin",
		"password": "admin",
	}

	t.Run("field match mutate", func(t *testing.T) {
		m := NewFieldMatchMutator()
		pm := &PasswordDefaultMutator{}
		m.Hook().Add("Password", pm)
		m.Hook().Add("password", pm)

		m.Execute(&cc)
		assert.Equal(t, cc.Username, "admin")
		assert.Equal(t, cc.Password, "********")

		m.Execute(&cm)
		assert.Equal(t, cm["username"], "admin")
		assert.Equal(t, cm["password"], "********")
	})
}
