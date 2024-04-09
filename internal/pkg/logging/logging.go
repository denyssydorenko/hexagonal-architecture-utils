package logging

import (
	"context"
	"fmt"
	domain "hexagonal-architexture-utils/internal/domains"
	"os"
	"time"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerWithCtx struct {
	*zap.Logger
}

var Global = createLogger()

func createLogger() LoggerWithCtx {
	stdout := zapcore.AddSync(os.Stdout)

	level := zap.NewAtomicLevelAt(zap.InfoLevel)

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	productionCfg.EncodeLevel = zapcore.CapitalLevelEncoder
	productionCfg.StacktraceKey = "stack"

	jsonEncoder := zapcore.NewJSONEncoder(productionCfg)

	jsonOutCore := zapcore.NewCore(jsonEncoder, stdout, level)

	samplingCore := zapcore.NewSamplerWithOptions(
		jsonOutCore,
		time.Second, // interval
		3,           // log first 3 entries
		0,           // thereafter log zero entires within the interval
	)

	return LoggerWithCtx{zap.New(samplingCore)}
}

func (l *LoggerWithCtx) logFields(ctx context.Context, fields []zap.Field) []zap.Field {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		context := span.SpanContext()
		spanField := zap.String("span_id", context.SpanID().String())
		traceField := zap.String("trace_id", context.TraceID().String())
		traceFlags := zap.Int("trace_flags", int(context.TraceFlags()))
		fields = append(fields, []zap.Field{spanField, traceField, traceFlags}...)
	}

	xtraceidString := fmt.Sprintf("%s", ctx.Value(domain.XTRACEID))
	fields = append(fields, []zap.Field{zap.String(domain.XTRACEID, xtraceidString)}...)

	return fields
}

func (log *LoggerWithCtx) Info(ctx context.Context, msg string, fields ...zap.Field) {
	fieldsWithTraceCtx := log.logFields(ctx, fields)
	log.Logger.Info(msg, fieldsWithTraceCtx...)
}

func (log *LoggerWithCtx) Error(ctx context.Context, msg string, fields ...zap.Field) {
	fieldsWithTraceCtx := log.logFields(ctx, fields)
	log.Logger.Error(msg, fieldsWithTraceCtx...)
}
