package mstracer

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"music-snap/pkg/app"
	"os"
	"time"
)

func Init(cfg *Config, appCfg *app.Config) (*sdktrace.TracerProvider, func() error, error) {
	ctx := context.Background()
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(appCfg.Name),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	var batcher sdktrace.TracerProviderOption
	if cfg.Enable {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		conn, err := grpc.DialContext(ctx,
			cfg.ExpTarget,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
		}

		exp, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create trace exporter: %w", err)
		}
		batcher = sdktrace.WithBatcher(exp)
	} else {
		exp, err := stdout.New(stdout.WithPrettyPrint(), stdout.WithWriter(log.Writer()))
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create trace exporter: %w", err)
		}
		if !cfg.StdOut {
			// write to cfg.TraceFile
			if cfg.TraceFile == "" {
				return nil, nil, fmt.Errorf("trace file not specified: %s", cfg.TraceFile)
			}
			file, err := os.Create(cfg.TraceFile)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to create trace file: %w", err)
			}

			fileClose := func() error {
				return file.Close()
			}

			exp, err = stdout.New(stdout.WithWriter(file))
			if err != nil {
				return nil, fileClose, fmt.Errorf("failed to create trace exporter fo trace file: %w", err)
			}
		}
		batcher = sdktrace.WithBatcher(exp)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		batcher,
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil, nil
}
