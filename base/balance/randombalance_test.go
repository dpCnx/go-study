package balance

import "testing"

func TestRandomBalance_Get(t *testing.T) {
	r := RandomBalance{}
	_ = r.Add("1")
	_ = r.Add("2")
	_ = r.Add("3")
	_ = r.Add("4")

	for i := 0; i <= 10; i++ {
		s := r.Get()
		t.Log(s)
	}

}
