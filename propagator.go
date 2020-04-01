package prop

import (
	"context"
	"fmt"
	"strings"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/propagation"
	"go.opentelemetry.io/otel/api/trace"
)

const (
	DefaultTraceContextHeader = "uber-trace-id"

	// DefaultJaegerDebugHeader  = "jaeger-debug-id"
	// TraceBaggageHeaderPrefix = "uberctx-"
	// JaegerBaggageHeader      = "jaeger-baggage"
)

var (
	DefaultJaeger = Jaeger{
		TraceContextHeader: DefaultTraceContextHeader,
	}
)

type Jaeger struct {
	TraceContextHeader string
}

func (j *Jaeger) Inject(ctx context.Context, supplier propagation.HTTPSupplier) {
	sc := trace.SpanFromContext(ctx).SpanContext()
	if !sc.IsValid() {
		return
	}

	var flag uint8
	if sc.IsSampled() {
		flag |= 0x01
	}

	h := fmt.Sprintf(
		"%s:%s:0:%d",
		sc.TraceIDString(),
		sc.SpanIDString(),
		flag,
	)
	supplier.Set(j.TraceContextHeader, h)
}

func (j Jaeger) Extract(ctx context.Context, supplier propagation.HTTPSupplier) context.Context {
	h := supplier.Get(j.TraceContextHeader)
	if h == "" {
		return trace.ContextWithRemoteSpanContext(ctx, core.EmptySpanContext())
	}

	parts := strings.Split(h, ":")
	if len(parts) != 4 {
		return trace.ContextWithRemoteSpanContext(ctx, core.EmptySpanContext())
	}

	var err error
	sc := core.SpanContext{}
	sc.TraceID, err = core.TraceIDFromHex(parts[0])
	if err != nil {
		return trace.ContextWithRemoteSpanContext(ctx, core.EmptySpanContext())
	}
	sc.SpanID, err = core.SpanIDFromHex(parts[1])
	if err != nil {
		return trace.ContextWithRemoteSpanContext(ctx, core.EmptySpanContext())
	}
	var flag, mask byte = 0, 0x01
	_, err = fmt.Sscanf(parts[3], "%d", &flag)
	if err != nil {
		return trace.ContextWithRemoteSpanContext(ctx, core.EmptySpanContext())
	}
	if flag&mask == 0x01 {
		sc.TraceFlags = core.TraceFlagsSampled
	}

	return trace.ContextWithRemoteSpanContext(ctx, sc)
}

func (j Jaeger) GetAllKeys() []string {
	return []string{j.TraceContextHeader}
}
