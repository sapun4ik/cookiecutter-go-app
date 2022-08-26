package tracing

import (
	"context"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func StartHTTPServerTracerSpan(c *gin.Context, operationName string) (context.Context, opentracing.Span) {
	spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
	if err != nil {
		serverSpan := opentracing.GlobalTracer().StartSpan(operationName)
		ctx := opentracing.ContextWithSpan(c.Request.Context(), serverSpan)
		return ctx, serverSpan
	}

	serverSpan := opentracing.GlobalTracer().StartSpan(operationName, ext.RPCServerOption(spanCtx))
	ctx := opentracing.ContextWithSpan(c.Request.Context(), serverSpan)

	return ctx, serverSpan
}

func InjectTextMapCarrier(spanCtx opentracing.SpanContext) (opentracing.TextMapCarrier, error) {
	m := make(opentracing.TextMapCarrier)
	if err := opentracing.GlobalTracer().Inject(spanCtx, opentracing.TextMap, m); err != nil {
		return nil, err
	}
	return m, nil
}

func ExtractTextMapCarrier(spanCtx opentracing.SpanContext) opentracing.TextMapCarrier {
	textMapCarrier, err := InjectTextMapCarrier(spanCtx)
	if err != nil {
		return make(opentracing.TextMapCarrier)
	}
	return textMapCarrier
}

func ExtractTextMapCarrierHeaders(spanCtx opentracing.SpanContext) map[string]string {
	textMapCarrier, err := InjectTextMapCarrier(spanCtx)
	if err != nil {
		return make(opentracing.TextMapCarrier)
	}
	return textMapCarrier
}

func ExtractTextMapCarrierBytes(spanCtx opentracing.SpanContext) []byte {
	textMapCarrier, err := InjectTextMapCarrier(spanCtx)
	if err != nil {
		return []byte("")
	}

	dataBytes, err := json.Marshal(&textMapCarrier)
	if err != nil {
		return []byte("")
	}
	return dataBytes
}

func TraceErr(span opentracing.Span, err error) {
	span.SetTag("error", true)
	span.LogKV("error_code", err.Error())
}

func TraceWithErr(span opentracing.Span, err error) error {
	if err != nil {
		span.SetTag("error", true)
		span.LogKV("error_code", err.Error())
	}

	return err
}
