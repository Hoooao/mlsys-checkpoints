package main

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
)

// Traced main
// TracerProvider link exporter and data collection
func main() {
	fmt.Println("Service opened")
	l := log.New(os.Stdout, "", 0)

	// Write telemetry data to a file.
	f, err := os.Create("traces.txt")
	if err != nil {
		l.Fatal(err)
	}
	defer f.Close()

	// new exporter
	// a console exporter that will export to a file
	exp, err := newExporter(f)
	if err != nil {
		l.Fatal(err)
	}

	// new Provider
	// registering the exporter with a new TracerProvider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(newResource()),
		trace.WithSampler(trace.NeverSample()),
	)
	sigCh := make(chan string, 1)
	//signal.Notify(sigCh, os.Interrupt)

	errCh := make(chan error)

	// deferring a function to flush and stop TP
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			l.Fatal(err)
		}
	}()

	go func() {
		http.HandleFunc("/", fibbo) // set routes
		http.HandleFunc("/finish", func(writer http.ResponseWriter, request *http.Request) {
			sigCh <- "terminate"
			return
		})
		http.HandleFunc("/finish-tracing", func(writer http.ResponseWriter, request *http.Request) {
			if err := tp.Shutdown(context.Background()); err != nil {
				l.Fatal(err)
			}
			return
		})
		fmt.Println("Service started")
		err2 := http.ListenAndServe(":9090", nil) // set port
		if err != nil {
			log.Fatal("ListenAndServe: ", err2)
		}
	}()
	// registering it as the global OpenTelemetry TracerProvider
	otel.SetTracerProvider(tp)

	select {
	case <-sigCh:
		l.Println("\ngoodbye")
		return
	case err := <-errCh:
		if err != nil {
			l.Fatal(err)
		}
	}

}

func fibbo(w http.ResponseWriter, r *http.Request) {
	// 原版运行-------------

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	errCh := make(chan error)
	app := NewApp(&w, r)
	go func() {
		errCh <- app.Run(context.Background())
	}()

}

/**
old

func main() {
	l := log.New(os.Stdout, "", 0)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	errCh := make(chan error)
	app := NewApp(os.Stdin, l)
	go func() {
		errCh <- app.Run(context.Background())
	}()

	select {
	case <-sigCh:
		l.Println("\ngoodbye")
		return
	case err := <-errCh:
		if err != nil {
			l.Fatal(err)
		}
	}
}
*/

// newResource returns a resource describing this application.
// 一个Resource represent the entity producing telemetry
func newResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("fib"),
			semconv.ServiceVersionKey.String("v0.1.0"),
			attribute.String("environment", "demo"),
		),
	)
	return r
}

// newExporter returns a console exporter.
func newExporter(w io.Writer) (trace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		// Use human-readable output.
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.
		//stdouttrace.WithoutTimestamps(),
	)
}
