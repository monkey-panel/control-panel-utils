package types

import (
	"sync"
	"testing"
)

func TestSummonID(t *testing.T) {
	s := NewSummonID()

	x := &sync.Map{}
	for i := 0; i < 200; i++ {
		go func() {
			y := s.Generate()
			str := y.Base2()

			t.Logf("\033[31m%v\033[0m-\033[34m%v\033[33m%v\033[0m", y.String(), str[:42], str[42:])
			if old, ok := x.Load(str); ok {
				t.Errorf("x(%d) & y(%d) are the same", old, y)
			}

			x.Store(str, y)
		}()
	}
}
