package main

import (
	"fmt"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/logging"
	"log"
)

func main() {

	conf := config.NewDefaultConfig()
	conf.Sentinel.Log.Logger = logging.NewConsoleLogger()
	err := sentinel.InitWithConfig(conf)
	if err != nil {
		log.Fatal(err)
	}

	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "some-test",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              2,
			StatIntervalInMs:       1000,
		},
	})
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err)
		return
	}

	for i := 0; i < 10; i++ {
		MyApi(i)
	}

}

func MyApi(i int) {
	e, b := sentinel.Entry("some-test", sentinel.WithTrafficType(base.Inbound))
	if b != nil {
		fmt.Println("no==>", i, b)
		// 请求被流控，可以从 BlockError 中获取限流详情
		// block 后不需要进行 Exit()
	} else {
		fmt.Println("come==>", i)
		// 请求可以通过，在此处编写您的业务逻辑
		// 务必保证业务逻辑结束后 Exit
		e.Exit()
	}
}
