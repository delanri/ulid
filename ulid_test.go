package ulid

import (
	"math"
	"reflect"
	"testing"
)

func TestUIDGen(t *testing.T) {
	guid := New(0, VARCHAR26)
	uid := guid.UID()

	p, ok := guid.Parse(uid.String())
	if !ok {
		t.Errorf("Generated UID %s failed to parse", uid)
	}
	if !reflect.DeepEqual(uid, p) {
		t.Errorf("Generated UID %s not equal with parse", uid)
	}

	exists := make(map[string]bool)
	for i := 0; i < int(math.Pow10(6)); i++ {
		uid := guid.UID().String()
		if exists[uid] {
			t.Errorf("Generated UID %s clashes at iteration %d", uid, i+1)
		}
		exists[uid] = true
	}
}