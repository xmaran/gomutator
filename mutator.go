// package gomutator provides a callback to the mutate hook function to mutate
// the value of the struct field or map key. This package would only work on
// addressable values. For struct types, it would work only on exposed struct
// fields
package gomutator

import (
	"fmt"
	"reflect"
	"sync"
)

type (
	// MutateHook provides contracts to mutate hook types. The consumer has to
	// write their own Mutate method to return new value to the matched field
	MutateHook interface {
		Mutate(any, any) any
	}

	MutatorChain struct {
		lock         sync.Mutex
		mutateMapper map[any]MutateHook
	}

	MutateType int8

	// Mutator provides contracts to mutate functionality
	Mutator interface {
		Execute(any)
		Hook() *MutatorChain
	}

	Mutate struct {
		hook      *MutatorChain
		matchType MutateType
	}
)

// String set
type strset map[string]struct{}

// Different match type mutator
const (
	TypeAndFieldMatch MutateType = iota + 1
	FieldMatch
)

// NewFieldMatchMutator creates a pointer to the Mutate type. Moreover, it
// initializes with the field name match mutate type and returns the Mutator
// interface
func NewFieldMatchMutator() Mutator {
	return NewMutator(FieldMatch)
}

// NewTypeAndFieldMatchMutator creates a pointer to the Mutate type. Moreover,
// it initializes with the type and field name match mutate type and returns the
// Mutator interface
func NewTypeAndFieldMatchMutator() Mutator {
	return NewMutator(TypeAndFieldMatch)
}

// NewMutator creates a pointer to the Mutate type. Moreover, it initializes
// with the input mutate type and returns the Mutator interface
func NewMutator(typ MutateType) Mutator {
	return &Mutate{
		matchType: typ,
		hook: &MutatorChain{
			mutateMapper: make(map[any]MutateHook),
		},
	}
}

func (mc *MutatorChain) Get(key any) (MutateHook, bool) {
	i, ok := mc.mutateMapper[key]
	return i, ok
}

func (mc *MutatorChain) Add(key any, mh MutateHook) *MutatorChain {
	mc.lock.Lock()
	defer mc.lock.Unlock()

	mc.mutateMapper[key] = mh
	return mc
}

func (mc *MutatorChain) Remove(key string) *MutatorChain {
	mc.lock.Lock()
	defer mc.lock.Unlock()

	delete(mc.mutateMapper, key)
	return mc
}

func (m *Mutate) Hook() *MutatorChain {
	return m.hook
}

func (m *Mutate) MatchType() MutateType {
	return m.matchType
}

// Execute execute the mutate hook function against the matched struct field name or map key
func (m *Mutate) Execute(data any) {

	// To prevent pointer loops while visiting addresses, a new set needs to
	// be created. This store will keep track of the visited addresses and help
	// in identifying whether the address has already been visited or not.
	visited := make(strset)
	// Process the data recursively
	m.processRecursively(reflect.ValueOf(data), visited)
}

func (m *Mutate) processRecursively(v reflect.Value, visited strset) {
	// Check whether the value address is already visited or not
	if ok := alreadyVisited(v, visited); ok {
		return
	}

	switch v.Kind() {
	case reflect.Invalid:
		return
	case reflect.Struct:
		m.processStructKind(v, visited)
	case reflect.Map:
		m.processMapKind(v, visited)
	case reflect.Ptr, reflect.Interface:
		m.processRecursively(v.Elem(), visited)
	}
}

func (m *Mutate) processStructKind(v reflect.Value, visited strset) {
	// If the kind is not a struct, return silently
	if v.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !field.CanSet() {
			continue
		}

		// default: field name match mutation
		key := v.Type().Field(i).Name
		// But if the match type is a type and a field match, concatenate the
		// actual struct type with the field name (for example,
		// restv1.Smtp.Username).
		if m.MatchType() == TypeAndFieldMatch {
			key = fmt.Sprintf("%T.%v", v.Interface(), key)
		}

		// Check whether the hook map has the key or not
		if mh, ok := m.Hook().Get(key); ok {
			// Call the mutate function with the underlined struct type and
			// field type to get the new value, which should be set to the
			// field.
			newValue := mh.Mutate(v.Interface(), field.Interface())
			field.Set(reflect.ValueOf(newValue))
		} else {
			// Otherwise, process the field recursively
			m.processRecursively(field, visited)
		}
	}
}

func (m *Mutate) processMapKind(v reflect.Value, visited strset) {
	// If the kind is not a map, return silently
	if v.Kind() != reflect.Map {
		return
	}

	if !v.CanSet() {
		return
	}

	for _, key := range v.MapKeys() {
		value := v.MapIndex(key)
		//  Check whether the hook map has the key or not
		if mh, ok := m.Hook().Get(key.Interface()); ok {
			// If the key is present, call mutate function with the underlined
			// value field type to get the new value, which should be set to
			// the key.
			newValue := mh.Mutate(nil, value.Interface())
			v.SetMapIndex(key, reflect.ValueOf(newValue))
		} else {
			// Otherwise, process the field recursively
			m.processRecursively(value, visited)
		}
	}
}

// alreadyVisited prevents pointer loops while visiting addresses
// The visited set will keep track of the visited addresses and help
// in identifying whether the address has already been visited or not.
func alreadyVisited(v reflect.Value, visited strset) bool {
	if !v.CanAddr() {
		return false
	}

	addr := v.Addr().String()
	if _, ok := visited[addr]; ok {
		return true
	}

	visited[addr] = struct{}{}
	return false
}
