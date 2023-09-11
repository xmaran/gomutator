package gomutator

type PasswordDefaultMutator struct{}

func (m PasswordDefaultMutator) Mutate(parent, current any) any {
	if _, ok := current.(string); !ok {
		return current
	}

	return "********"
}
