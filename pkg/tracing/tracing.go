package tracing

import (
	"context"
	"os"

	"github.com/go-logr/zapr"
	"github.com/rancher/opni/pkg/logger"
	"go.opentelemetry.io/contrib/propagators/autoprop"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log = logger.New(logger.WithLogLevel(zapcore.InfoLevel)).Named("tracing")

func Configure(serviceName string) {
	res, err := resource.New(context.Background(), resource.WithAttributes(
		semconv.ServiceNameKey.String(serviceName),
	))
	if err != nil {
		log.With(zap.Error(err)).Error("failed to configure tracing")
		return
	}

	opts := []tracesdk.TracerProviderOption{
		tracesdk.WithResource(res),
	}

	switch os.Getenv("OTEL_TRACES_EXPORTER") {
	case "jaeger":
		log.Info("using jaeger exporter")
		exp, err := jaeger.New(jaeger.WithCollectorEndpoint())
		if err != nil {
			log.With(zap.Error(err)).Error("failed to create exporter")
			return
		}
		opts = append(opts, tracesdk.WithBatcher(exp))
	case "otlp":
		log.Info("using otel exporter")
		exporter, err := otlptracegrpc.New(context.Background())
		if err != nil {
			log.With(zap.Error(err)).Error("failed to create exporter")
			return
		}
		opts = append(opts, tracesdk.WithBatcher(exporter))
	default:
	}

	otel.SetTracerProvider(tracesdk.NewTracerProvider(opts...))
	otel.SetTextMapPropagator(autoprop.NewTextMapPropagator())
	otel.SetLogger(zapr.NewLogger(log.Desugar()))
}
