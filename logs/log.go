package logs

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var (
	zapLog   *zap.SugaredLogger // 简易版日志文件
	logLevel = zap.NewAtomicLevel()
)

type Level int8

//日志级别
const (
	DebugLevel Level = iota - 1

	InfoLevel

	WarnLevel

	ErrorLevel

	DPanicLevel

	PanicLevel

	FatalLevel
)

type ZLog struct {
}

// InitLog 初始化日志文件
func New(path string) (err error) {

	cores := make([]zapcore.Core, 0)

	setLevel(InfoLevel)

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	//输入日志到文件
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   path,
		MaxSize:    128, // MB
		LocalTime:  true,
		Compress:   true,
		MaxBackups: 8, // 最多保留 n 个备份
	})

	file := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		w,
		logLevel,
	)
	cores = append(cores, file)
	//输入日志到控制台
	console := zapcore.NewCore(
		zapcore.NewConsoleEncoder(config),
		zapcore.AddSync(os.Stdout),
		logLevel,
	)
	cores = append(cores, console)
	logger := zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddCallerSkip(1))
	zapLog = logger.Sugar()

	return nil
}

func setLevel(level Level) {
	logLevel.SetLevel(zapcore.Level(level))
}

func Info(args ...interface{}) {
	zapLog.Info(args...)
}

func Infof(template string, args ...interface{}) {
	zapLog.Infof(template, args...)
}

func Warn(args ...interface{}) {
	zapLog.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	zapLog.Warnf(template, args...)
}

func Error(args ...interface{}) {
	zapLog.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	zapLog.Errorf(template, args...)
}
