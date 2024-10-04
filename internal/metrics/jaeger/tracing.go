package jaeger

import (
	"context"
	"github.com/JMURv/seo-svc/pkg/config"
	"github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
	"log"
)

func Start(ctx context.Context, serviceName string, conf *config.JaegerConfig) {
	cfg := jaeger.Configuration{
		ServiceName: serviceName,
		Sampler: &jaeger.SamplerConfig{
			Type:  conf.Sampler.Type,
			Param: float64(conf.Sampler.Param),
		},
		Reporter: &jaeger.ReporterConfig{
			LogSpans:           conf.Reporter.LogSpans,
			LocalAgentHostPort: conf.Reporter.LocalAgentHostPort,
		},
	}

	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		log.Fatalf("Error initializing Jaeger tracer: %s", err.Error())
	}

	opentracing.SetGlobalTracer(tracer)

	zap.L().Debug("Jaeger has been started")
	<-ctx.Done()

	zap.L().Debug("Shutting down Jaeger")
	if err = closer.Close(); err != nil {
		zap.L().Debug("Error shutting down Jaeger", zap.Error(err))
	}
	zap.L().Debug("Jaeger has been stopped")
}
