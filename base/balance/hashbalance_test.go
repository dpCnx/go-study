package balance

import (
	"fmt"
	"testing"
)

func TestHashBalance_Get(t *testing.T) {

	rb := NewHashBalance(2, nil)

	_ = rb.Add("1")
	_ = rb.Add("2")
	_ = rb.Add("3")
	_ = rb.Add("4")
	_ = rb.Add("5")

	for k, v := range rb.hashMap {
		t.Log(fmt.Sprintf("%v===>%v", k, v))
	}

	res, _ := rb.Get("8")
	t.Log(res)

	res2, _ := rb.Get("8")
	t.Log(res2)
}
