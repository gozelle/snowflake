package snowflake

import (
	"testing"
)

func TestDefaultGenerator(t *testing.T) {
	gen, err := DefaultGenerator(1)
	if err != nil {
		t.Error(err)
		return
	}

	id, e := gen.Generate()
	if e != nil {
		t.Error(e)
		return
	}
	t.Log(gen.EndAt())
	t.Log(id)
}
