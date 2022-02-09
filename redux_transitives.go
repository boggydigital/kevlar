package kvas

import "fmt"

//reduxTransitives are connections between values, similar to joins, that happen automatically.
//Consider an object that has field `language`, where the value is one of the language codes (`en`, `ru`, ...).
//Consider a map that provides human-readable name of the language by code (`en`: `English`, `ru`: `Русский`, ...).
//This logical connection can be established by setting those fields as transitive.
//When that happens, clients will transparently get `Language name (code)` for a `language` field.
type reduxTransitives map[string]string

func (rt reduxTransitives) IsTransitive(key string) bool {
	for t, _ := range rt {
		if key == t {
			return true
		}
	}
	return false
}

func (rt reduxTransitives) Transition(key string) string {
	return rt[key]
}

func (rt reduxTransitives) Fmt(from, to string) string {
	return fmt.Sprintf("%s (%s)", to, from)
}
