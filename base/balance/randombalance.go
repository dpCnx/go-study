package balance

import (
	"errors"
	"math/rand"
)

// RandomBalance 随机负载均衡
type RandomBalance struct {
	curIndex int

	rss []string
}

func (r *RandomBalance) Add(url ...string) error {
	if len(url) <= 0 {
		return errors.New("请添加url")
	}
	r.rss = append(r.rss, url...)

	return nil
}

func (r *RandomBalance) Get() string {
	return r.next()
}

func (r *RandomBalance) next() string {
	if len(r.rss) == 0 {
		return ""
	}
	r.curIndex = rand.Intn(len(r.rss))
	return r.rss[r.curIndex]
}
