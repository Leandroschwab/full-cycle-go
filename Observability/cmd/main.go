package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"time"

	"os"

	"github.com/Leandroschwab/full-cycle-go/Observability/internal/handlers"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"

	//"go.opentelemetry.io/otel/sdk"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	senconv "go.opentelemetry.io/otel/semconv/v1.21.0"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func initProvider(serviceName, collectorURL string) (func(context.Context) error, error) {
	ctx := context.Background()
	res, err := sdkresource.New(ctx,
		sdkresource.WithAttributes(
			senconv.ServiceNameKey.String(serviceName),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	/*conn, err := grpc.DialContext(ctx, collectorURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)*/

	conn, err := grpc.NewClient(collectorURL, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	bsb := sdktrace.NewBatchSpanProcessor(traceExporter)

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsb),
	)

	otel.SetTracerProvider(tracerProvider)

	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tracerProvider.Shutdown, nil
}

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	shutdown, err := initProvider(os.Getenv("OTEL_SERVICE_NAME"), os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
	if err != nil {
		log.Fatalf("Failed to shutdown TracerProvider: %v\n", err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Printf("Error shutting down TracerProvider: %v\n", err)
		}
	}()

	service := os.Getenv("FUNCTION")
	http.HandleFunc("/temperature", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		switch service {
		case "orchestrator":
			handlers.HandleCEPCode(w, r)
		case "inputvalidator":
			handlers.ValidateCEPCode(w, r)
		default:
			http.Error(w, "Invalid function specified", http.StatusInternalServerError)
		}
	})

	// Criar o servidor explicitamente
	server := &http.Server{
		Addr:    ":" + os.Getenv("HTTP_PORT"),
		Handler: nil, // DefaultServeMux is used
	}

	// Iniciar o servidor em uma goroutine
	go func() {
		log.Println("Starting server on", os.Getenv("HTTP_PORT"))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %s\n", err)
		}
	}()

	// Aguardar sinal de interrupção
	select {
	case <-sigCh:
		log.Println("Shutting down gracefully, CTRL+C pressed")
	case <-ctx.Done():
		log.Println("Shutting down due to other reason")
	}

	// Criar contexto com timeout para desligamento gracioso
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*5)
	defer shutdownCancel()

	// Desligar o servidor adequadamente
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v\n", err)
	}

	log.Println("Server gracefully stopped")
}
