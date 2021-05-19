package middleware

import (
	"fmt"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/util"
	"github.com/gin-gonic/gin"
	"go-study/gin/models"
	"go.uber.org/zap"
)

func init() {

	conf := config.NewDefaultConfig()
	// conf.Sentinel.Log.Logger = logging.NewConsoleLogger()
	if err := sentinel.InitWithConfig(conf); err != nil {
		zap.L().Error(err.Error())
		return
	}

	_, err := flow.LoadRules([]*flow.Rule{
		{
			Resource:               "web",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              1000,
			StatIntervalInMs:       1000, // 表示统计周期是1s
		},
	})

	if err != nil {
		zap.L().Error(err.Error())
		return
	}

	circuitbreaker.RegisterStateChangeListeners(&stateChangeTestListener{})

	_, err = circuitbreaker.LoadRules([]*circuitbreaker.Rule{
		{
			Resource:         "web",
			Strategy:         circuitbreaker.SlowRequestRatio,
			RetryTimeoutMs:   3000, // 熔断触发后持续的时间
			MinRequestAmount: 3,
			StatIntervalMs:   1000, // 统计的时间窗口长度（单位为 ms）
			MaxAllowedRtMs:   50,   // 如果response time大于MaxAllowedRtMs，那么当前请求就属于慢调用。
			Threshold:        0.1,
		},
	})

	if err != nil {
		zap.L().Error(err.Error())
		return
	}
}

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

func Sentinel() gin.HandlerFunc {

	return func(c *gin.Context) {

		e, b := sentinel.Entry("web", sentinel.WithTrafficType(base.Inbound))
		if b != nil {
			// 请求被流控，可以从 BlockError 中获取限流详情
			// block 后不需要进行 Exit()
			models.ResponseErrorWithMsg(c, models.CodeServerBusy, b.Error())
			c.Abort()
			return
		} else {
			// 请求可以通过，在此处编写您的业务逻辑
			c.Next()
			// 务必保证业务逻辑结束后 Exit
			e.Exit()
		}

	}

}
