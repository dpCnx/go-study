package main

import (
	"fmt"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/logging"
	"github.com/alibaba/sentinel-golang/util"
	"log"
	"time"
)

type stateChangeTestListener struct {
}

func (s *stateChangeTestListener) OnTransformToClosed(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	fmt.Printf("rule.steategy: %+v, From %s to Closed, time: %d\n", rule.Strategy, prev.String(), util.CurrentTimeMillis())
}

func (s *stateChangeTestListener) OnTransformToOpen(prev circuitbreaker.State, rule circuitbreaker.Rule, snapshot interface{}) {
	fmt.Printf("rule.steategy: %+v, From %s to Open, snapshot: %.2f, time: %d\n", rule.Strategy, prev.String(), snapshot, util.CurrentTimeMillis())
}

func (s *stateChangeTestListener) OnTransformToHalfOpen(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	fmt.Printf("rule.steategy: %+v, From %s to Half-Open, time: %d\n", rule.Strategy, prev.String(), util.CurrentTimeMillis())
}

func main() {
	conf := config.NewDefaultConfig()
	conf.Sentinel.Log.Logger = logging.NewConsoleLogger()
	err := sentinel.InitWithConfig(conf)
	if err != nil {
		log.Fatal(err)
	}
	circuitbreaker.RegisterStateChangeListeners(&stateChangeTestListener{})

	_, err = circuitbreaker.LoadRules([]*circuitbreaker.Rule{
		{
			Resource:         "abc",
			Strategy:         circuitbreaker.SlowRequestRatio,
			RetryTimeoutMs:   3000, //熔断触发后持续的时间
			MinRequestAmount: 3,
			StatIntervalMs:   1000, //统计的时间窗口长度（单位为 ms）
			MaxAllowedRtMs:   50,   //如果response time大于MaxAllowedRtMs，那么当前请求就属于慢调用。
			Threshold:        0.1,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		MyApi2(i)
	}
}

func MyApi2(i int) {
	e, b := sentinel.Entry("abc", sentinel.WithTrafficType(base.Inbound))
	if b != nil {
		fmt.Println("no==>", i, b)
		// 请求被流控，可以从 BlockError 中获取限流详情
		// block 后不需要进行 Exit()
	} else {
		fmt.Println("come==>", i)
		// 请求可以通过，在此处编写您的业务逻辑
		// 务必保证业务逻辑结束后 Exit
		if i == 1 {
			time.Sleep(60 * time.Millisecond)
		}
		e.Exit()
	}
}
