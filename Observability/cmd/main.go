package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Leandroschwab/full-cycle-go/Observability/internal/handlers"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func initProvider(serviceName, collectorURL string) (func(context.Context) error, error) {
	fmt.Println("Setting up the OpenTelemetry provider...")
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(semconv.ServiceName(serviceName)),
	)
	if err != nil {
		return nil, fmt.Errorf("resource creation error: %w", err)
	}

	// Use an increased timeout for dialing the collector
	ctx, cancel := context.WithTimeout(ctx, 12*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, collectorURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial collector: %w", err)
	}

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	spProcessor := sdktrace.NewBatchSpanProcessor(traceExporter)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(spProcessor),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tp.Shutdown, nil
}

func init() {
	viper.AutomaticEnv()
}

func main() {
	log.Println("Microservice starting...")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	shutdown, err := initProvider(viper.GetString("OTEL_SERVICE_NAME"), viper.GetString("OTEL_EXPORTER_OTLP_ENDPOINT"))
	if err != nil {
		log.Fatalf("Provider init error: %v", err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatalf("Error shutting down TracerProvider: %v", err)
		}
	}()

	tracer := otel.Tracer("custom-temperature-service")
	templateData := &handlers.TemplateData{
		Funtion:    viper.GetString("FUNCTION"),
		HTTP_PORT:  viper.GetString("HTTP_PORT"),
		OTELTracer: tracer,
	}

	// Use our custom Handler (server) instance.
	server := handlers.NewServer(templateData)

	// Wrap the /temperature endpoint with otelhttp for automatic instrumentation.
	http.Handle("/temperature", otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch templateData.Funtion {
		case "orchestrator":
			handlers.HandleCEPCode(w, r)
		case "inputvalidator":
			// Now calling the method on our server instance.
			server.ValidateCEPCode(w, r)
		default:
			http.Error(w, "Invalid function specified", http.StatusInternalServerError)
		}
	}), "/temperature"))

	srvAddr := ":" + templateData.HTTP_PORT
	srv := &http.Server{
		Addr:    srvAddr,
		Handler: nil, // default mux
	}

	go func() {
		log.Printf("Server running on %s...", srvAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	select {
	case <-sigCh:
		log.Println("Interrupt signal received, shutting down...")
	case <-ctx.Done():
		log.Println("Context cancelled, shutting down...")
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Graceful shutdown error: %v", err)
	}

	log.Println("Server stopped gracefully")
}
