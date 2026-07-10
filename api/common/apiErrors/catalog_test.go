package apiErrors

import "testing"

func TestCodesAreUnique(t *testing.T) {
	seen := make(map[string]string, len(All))
	for _, e := range All {
		if e == nil {
			t.Fatalf("nil entry in All")
		}
		if prev, ok := seen[e.Code]; ok {
			t.Errorf("duplicate code %s: %q and %q", e.Code, prev, e.Message)
			continue
		}
		seen[e.Code] = e.Message
	}
}
