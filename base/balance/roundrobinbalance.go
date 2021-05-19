package balance

import (
	"errors"
)

// RoundRobinBalance 轮询负载均衡
type RoundRobinBalance struct {
	curIndex int
	rss      []string
}

func (r *RoundRobinBalance) Add(url ...string) error {

	if len(url) <= 0 {
		return errors.New("请添加url")
	}
	r.rss = append(r.rss, url...)

	return nil
}

func (r *RoundRobinBalance) Get() string {
	return r.next()
}

func (r *RoundRobinBalance) GetAll() []string {
	return r.rss
}

func (r *RoundRobinBalance) Reset() {
	r.rss = r.rss[:0]
}

func (r *RoundRobinBalance) next() string {

	if len(r.rss) == 0 {
		return ""
	}
	lens := len(r.rss)
	if r.curIndex >= lens {
		r.curIndex = 0
	}

	curAddr := r.rss[r.curIndex]
	r.curIndex = (r.curIndex + 1) % lens
	return curAddr
}
