package middleware

import (
	"context"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"go-study/gin/conf"
)

/*
	const，全量采集。param采样率设置0,1 分别对应打开和关闭
	probabilistic ，概率采集。param默认万份之一，0~1之间取值，
	rateLimiting ，限速采集。param每秒采样的个数
	remote 动态采集策略。param值于probabilistic的参数一样。在收到实际值之前的初始采样率。改值可以通过环境变量的JAEGER_SAMPLER_PARAM设定
*/

func NewJaegerTracer() (opentracing.Tracer, io.Closer) {

	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const", // 固定采样
			Param: 1,       // 1=全采样、0=不采样
		},

		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: fmt.Sprintf("127.0.0.1:6831"),
		},

		ServiceName: conf.C.Server.Name,
	}

	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(err)
	}
	opentracing.SetGlobalTracer(tracer)
	return tracer, closer
}

func Jaeger() gin.HandlerFunc {

	return func(c *gin.Context) {
		tracer, closer := NewJaegerTracer()
		defer closer.Close()

		var newCtx context.Context
		var span opentracing.Span
		spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		if err != nil {
			span, newCtx = opentracing.StartSpanFromContextWithTracer(c.Request.Context(), tracer, c.Request.URL.Path)
		} else {
			span, newCtx = opentracing.StartSpanFromContextWithTracer(
				c.Request.Context(),
				tracer,
				c.Request.URL.Path,
				opentracing.ChildOf(spanCtx),
				opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
			)
		}
		defer span.Finish()

		var traceID string
		var spanID string
		var spanContext = span.Context()
		switch spanContext.(type) {
		case jaeger.SpanContext:
			jaegerContext := spanContext.(jaeger.SpanContext)
			traceID = jaegerContext.TraceID().String()
			spanID = jaegerContext.SpanID().String()
		}
		c.Set("X-Trace-ID", traceID)
		c.Set("X-Span-ID", spanID)
		c.Request = c.Request.WithContext(newCtx)
		c.Next()
	}
}
