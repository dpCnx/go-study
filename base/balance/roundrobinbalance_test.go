package balance

import "testing"

func TestRoundRobinBalance_Get(t *testing.T) {

	r := RoundRobinBalance{}

	_ = r.Add("127.0.0.1:9988", "127.0.0.1:9998")

	for i := 0; i < 10; i++ {
		t.Log(r.curIndex)
		t.Log(r.Get())
	}
}
