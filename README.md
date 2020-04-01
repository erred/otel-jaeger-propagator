# otel-jaeger-propagator

Jaeger propagator for OpenTelemetry

[![License](https://img.shields.io/github/license/seankhliao/otel-jaeger-propagator.svg?style=flat-square)](LICENSE)
[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/seankhliao/otel-jaeger-propagator)

## Usage

```go
package main

import (
        "net/http"

        prop "github.com/seankhliao/otel-jaeger-propagator"
        "go.opentelemetry.io/otel/api/global"
        "go.opentelemetry.io/otel/api/propagation"
)

func main() {
        // global propagation
        global.SetPropagators(
                propagation.WithExtractors(prop.DefaultJaeger),
                propagation.WithInjectors(prop.DefaultJaeger),
        )
}
```

## Propagation format

```txt
{trace-id}:{span-id}:{parent-span-id}:{flags}
```

Propagation is mapped to the [core.SpanContext](https://godoc.org/go.opentelemetry.io/otel/api/core#SpanContext) fields:

- `trace-id`: SpanContext.TraceID
- `span-oid`: SpanContext.SpanID
- `parent-span-id`: ignored, injected as `0`
- `flags`: SpanContext.TraceFlags
  - 0x01 is mapped to core.TraceFlagsSampled
  - all other flags are ignored

## notes

- [jaeger propagation format][jprop]

[jprop]: https://www.jaegertracing.io/docs/1.17/client-libraries/#propagation-format
