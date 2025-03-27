package logger

import (
	"fmt"
	configuration "github.com/csturiale/config-mgtm"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/url"

	"gopkg.in/natefinch/lumberjack.v2"
	logging "log"
	"os"
	"strings"
)

const logExtension = ".log"

var (
	log               *zap.SugaredLogger
	logFolder         = getString("log.folder")
	enableFileLogging = getBoolOrDefault("log.file.enable", true)
	filename          = getStringOrDefault("log.file.name", fmt.Sprintf("%s%s", "application", logExtension))
)

type lumberjackSink struct {
	*lumberjack.Logger
}

func (lumberjackSink) Sync() error {
	return nil
}

func init() {
	if !strings.HasSuffix(filename, logExtension) {
		filename = fmt.Sprintf("%s%s", filename, logExtension)
	}

	outputPaths := []string{"stdout"}

	errorPaths := []string{"stderr"}

	if enableFileLogging {
		err := os.MkdirAll(logFolder, os.ModePerm)
		if err != nil {
			logging.Fatalf("Error creating log folder: %v", err)
		}

		ll := lumberjack.Logger{
			Filename:   fmt.Sprintf("%s/%s", logFolder, filename),
			MaxSize:    configuration.GetIntOrDefault("log.file.maxSize", 10),     // Max size in MB
			MaxBackups: configuration.GetIntOrDefault("log.file.maxBackups", 3),   // Max number of old log files to keep
			MaxAge:     configuration.GetIntOrDefault("log.file.maxAge", 28),      // Max age in days to keep a log file
			Compress:   configuration.GetBoolOrDefault("log.file.compress", true), // Compress old log files
		}
		err = zap.RegisterSink("lumberjack", func(*url.URL) (zap.Sink, error) {
			return lumberjackSink{
				Logger: &ll,
			}, nil
		})
		if err != nil {
			panic(fmt.Sprintf("build zap logger from config error: %v", err))
		}
		outputPaths = append(outputPaths, fmt.Sprintf("lumberjack:%s", fmt.Sprintf("%s/%s", logFolder, filename)))
	}

	cfg := zap.Config{
		Encoding:         "console",
		Level:            zap.NewAtomicLevelAt(getLogLevel()),
		OutputPaths:      outputPaths,
		ErrorOutputPaths: errorPaths,
		Development:      getBoolOrDefault("log.dev", false),
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:      "level",
			TimeKey:       "time",
			EncodeTime:    zapcore.ISO8601TimeEncoder,
			CallerKey:     "caller",
			EncodeCaller:  zapcore.ShortCallerEncoder,
			NameKey:       "logger",
			MessageKey:    "msg",
			StacktraceKey: "stacktrace",
			LineEnding:    zapcore.DefaultLineEnding,
			//EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
		},
	}
	logger, err := cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(fmt.Sprintf("build zap logger from config error: %v", err))
	}
	defer logger.Sync() // flushes buffer, if any
	log = logger.Sugar()
}

// Info ...
func Info(v ...interface{}) {
	log.Infoln(v...)
}

// Warn ...
func Warn(v ...interface{}) {
	log.Warnln(v...)
}

// Debug ...
func Debug(v ...interface{}) {
	log.Debugln(v...)
}

// Trace ...
func Trace(v ...interface{}) {
	log.Debugln(v...)
}

// Error ...
func Error(v ...interface{}) {
	log.Errorln(v...)
}

// Fatal ...
func Fatal(v ...interface{}) {
	log.Fatalln(v...)
}

// Infof ...
func Infof(format string, v ...interface{}) {
	log.Infof(format, v...)
}

// Warnf ...
func Warnf(format string, v ...interface{}) {
	log.Warnf(format, v...)
}

// Debugf ...
func Debugf(format string, v ...interface{}) {
	log.Debugf(format, v...)
}

// Tracef ...
func Tracef(format string, v ...interface{}) {
	log.Debugf(format, v...)
}

// Errorf ...
func Errorf(format string, v ...interface{}) {
	log.Errorf(format, v...)
}

// Fatalf ...
func Fatalf(format string, v ...interface{}) {
	//fields := getFields(runtime.Caller(1))
	//log.WithFields(fields).Fatalf(format, v...)
	log.Fatalf(format, v...)
}

func GetInstance() *zap.SugaredLogger {
	return log
}

func getString(key string) string {
	//Debugf("Returning item %s", key)
	return viper.GetString(key)
}

func getStringOrDefault(key string, val string) string {
	//Debugf("Returning item %s", key)
	if viper.Get(key) == nil {
		return val
	}
	return viper.GetString(key)
}

func getBoolOrDefault(key string, value bool) bool {
	if viper.Get(key) == nil {
		return value
	}
	return viper.GetBool(key)
}

func getLogLevel() zapcore.Level {
	lvl := strings.ToLower(getStringOrDefault("log.level", "info"))
	switch lvl {
	case "info":
		return zapcore.InfoLevel
	case "debug":
		return zapcore.DebugLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel

	}
}
