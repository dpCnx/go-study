package balance

import (
	"testing"
)

func TestLB(t *testing.T) {

	rb := &WeightRoundRobinBalance{}
	_ = rb.Add("1", 3)
	_ = rb.Add("2", 2)
	_ = rb.Add("3", 1)

	for i := 0; i < 10; i++ {
		t.Log(rb.Get())
	}
}
