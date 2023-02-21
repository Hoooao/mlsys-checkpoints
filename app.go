package main

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"strconv"
)

// name is the Tracer name used to identify this instrumentation library.
const name = "fib"

// App is a Fibonacci computation application.
type App struct {
	w *http.ResponseWriter
	r *http.Request
}

// NewApp returns a new App.
func NewApp(w *http.ResponseWriter, r *http.Request) *App {
	return &App{w: w, r: r}
}

// Run starts polling users for Fibonacci number requests and writes results.
// Traced版本
// Tracer相当于放在了一个"服务始发点"上，每次这个服务(for-loop)运行都是一个新的trace/root span
// context从main而来，可以理解为trace本体
func (a *App) Run(ctx context.Context) error {

	// Each execution of the run loop, we should get a new "root" span and context.
	newCtx, span := otel.Tracer(name).Start(ctx, "Run")

	n, err := a.GetInput(newCtx)
	if err != nil {
		span.End()
		return err
	}

	a.Write(newCtx, n)
	span.End()
	return nil
}

/**
// 原版
// Run starts polling users for Fibonacci number requests and writes results.
func (a *App) Run(ctx context.Context) error {
	for {
		n, err := a.Poll(ctx)
		if err != nil {
			return err
		}

		a.Write(ctx, n)
	}
}
*/

// GetInput 得到用户输入
// Traced 版本
func (a *App) GetInput(ctx context.Context) (uint64, error) {
	// 再次新增一个span
	_, span := otel.Tracer(name).Start(ctx, "GetInput")
	defer span.End()
	a.r.ParseForm()

	n, err := strconv.Atoi(a.r.FormValue("n"))

	// Store n as a string to not overflow an int64.
	nStr := strconv.FormatUint(uint64(n), 10)
	span.SetAttributes(attribute.String("request.n", nStr))

	return uint64(n), err
}

/** 原版
func (a *App) Poll(ctx context.Context) (uint, error) {
	a.l.Print("What Fibonacci number would you like to know: ")

	var n uint
	_, err := fmt.Fscanf(a.r, "%d\n", &n)
	return n, err
}

*/

// Write writes the n-th Fibonacci number back to the user.
// Traced版本
// 这里新建了两个span：Write和Fib的
func (a *App) Write(ctx context.Context, n uint64) {
	var span trace.Span
	ctx, span = otel.Tracer(name).Start(ctx, "Write")
	defer span.End()

	f, err := func(ctx context.Context) (uint64, error) {
		_, span := otel.Tracer("12").Start(ctx, "Fibonacci")
		defer span.End()
		return Fibonacci(uint(n))
	}(ctx)
	if err != nil {
		fmt.Print("Fibonacci(%d): %v\n", n, err)
		fmt.Fprintf(*(a.w), "hello\n")
	} else {
		fmt.Printf("Fibonacci(%d) = %d\n", n, f)
		fmt.Fprintf(*(a.w), "hello\n")
	}
}

/**原版
func (a *App) Write(ctx context.Context, n uint) {
	f, err := Fibonacci(n)
	if err != nil {
		a.l.Printf("Fibonacci(%d): %v\n", n, err)
	} else {
		a.l.Printf("Fibonacci(%d) = %d\n", n, f)
	}
}
*/

// Context中的Tree结构：
/**
Run
├── GetInput
└── Write
    └── Fibonacci
*/
