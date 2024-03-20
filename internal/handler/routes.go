package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"dreampicai/cmd/web"
	"dreampicai/cmd/web/view/home"

	_ "dreampicai/cmd/web"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewMux()
	ctx := context.Background()

	exporter, err := otlptrace.New(
		ctx,
		otlptracehttp.NewClient(),
	)
	if err != nil {
		log.Fatalf("failed to create exporter: %v", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("dreampicai"),
		)),
	)
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatalf("error shutting down tracer provider: %v", err)
		}
	}()
	otel.SetTracerProvider(tp)
	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/health"))
	r.Use(middleware.Recoverer)
	r.Use(WithUser)
	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.FS(web.Files))))

	r.Get("/", s.HelloWorldHandler)

	r.Get("/health", s.healthHandler)
	r.Get("/home", MakeHandler("home", s.HandleHomeIndex))
	r.Get("/login", MakeHandler("login", s.HandleLoginIndex))
	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}

func (s *Server) HandleHomeIndex(w http.ResponseWriter, r *http.Request) error {
	return home.Index().Render(r.Context(), w)
}
