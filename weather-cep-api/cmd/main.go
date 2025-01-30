package main

import (
	"context"
	"log"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/joho/godotenv"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/weather-cep-api/internal/adapters"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/weather-cep-api/internal/handlers"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error trying to load env variables")
		return
	}

	tp, err := initTracer()
	if err != nil {
		log.Printf("Failed to initialize tracer: %v. Proceeding without tracing.", err)
	}

	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	adapter := &adapters.DefaultWeatherAdapter{}

	http.HandleFunc("/temperature/", handlers.TemperatureHandler(adapter))
	http.HandleFunc("/docs", handlers.DocsHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs", http.StatusFound)
	})

	log.Printf("Server running on port %s", "8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initTracer() (*sdktrace.TracerProvider, error) {
	exporter, err := zipkin.New("http://zipkin:9411/api/v2/spans")
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("weather-cep-api"),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp, nil
}
