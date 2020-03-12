package log

import (
	"context"
	"gitlab.intelligentb.com/devops/bplus/config"
	bplusc "gitlab.intelligentb.com/devops/bplus/context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger
func init(){
	initZap()
}

func initZap(){
	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(determineLogLevel()),
		DisableCaller:    getDisableCaller(),
		DisableStacktrace: getDisableStackTrace(),
		Development:getDevelopmentMode(),
		Encoding:         "json",
		EncoderConfig:    encoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		InitialFields: map[string]interface{}{
			"application": config.GetApplicationName(),
		},
	}

	var err error
	logger, err = cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger.Info("logger construction succeeded")

}

func Sync(){
	logger.Sync()
}

func determineLogLevel() zapcore.Level{
	loglevel := config.Value("bplus.log_level")
	if loglevel == "" {
		return zapcore.InfoLevel
	}

	switch loglevel{
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	default:
		return zapcore.InfoLevel
	}
}

func getDisableCaller()bool {
	return config.BoolValue("bplus.disable_caller")
}

func getDisableStackTrace()bool {
	return config.BoolValue("bplus.disable_stacktrace")
}

func getDevelopmentMode()bool {
	return config.BoolValue("bplus.development_mode")
}


func encoderConfig() zapcore.EncoderConfig {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.MessageKey = "message"
	encoderConfig.CallerKey = "caller"
	encoderConfig.LevelKey = "level"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.NameKey = "name"
	return encoderConfig
}

func Info(ctx context.Context, message string){
	logger.Info(message,enhanceContext(ctx)...)
}

func Error(ctx context.Context, message string){
	logger.Error(message,enhanceContext(ctx)...)
}

func Warn(ctx context.Context, message string){
	logger.Warn(message,enhanceContext(ctx)...)
}

func InfoWithFields(ctx context.Context,fields map[string]string,  message string){
	logger.Info(message,enhance(ctx,fields)...)
}

func ErrorWithFields(ctx context.Context,fields map[string]string, message string){
	logger.Error(message,enhance(ctx,fields)...)
}

func WarnWithFields(ctx context.Context,fields map[string]string, message string,args ...interface{}){
	logger.Warn(message,enhance(ctx,fields)...)
}

func enhance(ctx context.Context,fields map[string]string) [] zap.Field{
	ret := enhanceContext(ctx)
	for n,v := range fields {
		ret = append(ret,zap.String(n,v))
	}
	return ret
}

func enhanceContext(ctx context.Context)[]zap.Field{
	return []zap.Field{
		zap.String("TraceID", bplusc.Value(ctx, bplusc.TraceID).(string)),
	}
}
