package _uuid

import (
	"testing"
)

func TestDefaultGenerator(t *testing.T) {
	gen, err := NewSnowflake(1)
	if err != nil {
		t.Error(err)
		return
	}
	
	id, e := gen.ID()
	if e != nil {
		t.Error(e)
		return
	}
	t.Log(gen.EndAt())
	t.Log(id)
}
