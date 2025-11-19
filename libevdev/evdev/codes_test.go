package evdev

import (
	"testing"
)

type evCodeToString = map[EvCode]string
type evPropToString = map[EvProp]string
type evTypeToString = map[EvType]string

type evCodeFromString = map[string]EvCode
type evPropFromString = map[string]EvProp
type evTypeFromString = map[string]EvType

type testPair struct {
	name       string
	toString   interface{}
	fromString interface{}
}

var testingSet = []testPair{
	{"INPUT", INPUTToString, INPUTFromString},
	{"EV", EVToString, EVFromString},
	{"SYN", SYNToString, SYNFromString},
	{"KEY", KEYToString, KEYFromString},
	{"REL", RELToString, RELFromString},
	{"ABS", ABSToString, ABSFromString},
	{"SW", SWToString, SWFromString},
	{"MSC", MSCToString, MSCFromString},
	{"LED", LEDToString, LEDFromString},
	{"REP", REPToString, REPFromString},
	{"SND", SNDToString, SNDFromString},
	{"ID", IDToString, IDFromString},
	{"BUS", BUSToString, BUSFromString},
	{"MT", MTToString, MTFromString},
	{"FF", FFToString, FFFromString},
}

var validMapping, invalidMapping, duplicationMapping int

func contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func compareGroups[V EvCode | EvType | EvProp](t *testing.T, from map[string]V, to map[V]string) {
	var duplicatedVal = make(map[V][]string)

	for str1, val1 := range from {
		for val2, str2 := range to {
			if val1 == val2 && str1 != str2 {
				if !contains(duplicatedVal[val1], str1) {
					duplicatedVal[val1] = append(duplicatedVal[val1], str1)
				}
				if !contains(duplicatedVal[val1], str2) {
					duplicatedVal[val1] = append(duplicatedVal[val1], str2)
				}
			}
		}
	}

	// test mapping in to>from direction
	for k1, v := range to {
		k2, ok := from[v]
		if !ok {
			t.Logf("- missing fromString item! key: %v (expected value: %4d / 0x%04x)", v, k1, k1)
			invalidMapping++
			continue
		}
		if k1 != k2 {
			t.Logf("- different fromString value! key: %v, expected: %v, got: %v", v, k1, k2)
			invalidMapping++
			continue
		}
		validMapping++
	}

	// test mapping in from>to direction
	for k1, v := range from {
		k2, ok := to[v]
		if !ok {
			t.Logf("- missing toString item! key: %4d / 0x%04x (expected value: %v)", v, v, k1)
			invalidMapping++
			continue
		}
		if k1 != k2 {
			// may be different due to value duplication
			if contains(duplicatedVal[v], k1) {
				duplicationMapping++
				continue
			}
			t.Logf("- different toString value! key: %v, expected: %v, got: %v", v, k1, k2)
			invalidMapping++
			continue
		}

		validMapping++
	}
}

func TestCodesMappings(t *testing.T) {
	for _, p := range testingSet {
		t.Logf("Analyzing \"%s\" group...", p.name)
		switch to := p.toString.(type) {
		case evCodeToString:
			switch from := p.fromString.(type) {
			case evCodeFromString:
				compareGroups(t, from, to)
			case evPropFromString, evTypeFromString:
				t.Fatal("type mismatch")
			default:
				t.Fatal("unexpected type")
			}
		case evPropToString:
			switch from := p.fromString.(type) {
			case evPropFromString:
				compareGroups(t, from, to)
			case evCodeFromString, evTypeFromString:
				t.Fatal("type mismatch")
			default:
				t.Fatal("unexpected type")
			}
		case evTypeToString:
			switch from := p.fromString.(type) {
			case evTypeFromString:
				compareGroups(t, from, to)
			case evCodeFromString, evPropFromString:
				t.Fatal("type mismatch")
			default:
				t.Fatal("unexpected type")
			}
		default:
			t.Fatal("unexpected type")
		}
	}

	t.Logf(
		"valid mappings: %d, ignored due to key value duplication: %d, invalid mappings: %d",
		validMapping, duplicationMapping, invalidMapping,
	)
	if invalidMapping > 0 {
		t.Fatal("detected invalid mappings!")
	}

	t.Log("all good")
}
