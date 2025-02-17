package logger

import (
	"fmt"
	configuration "github.com/csturiale/config-mgtm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	logging "log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const logExtension = ".log"

var (
	log               *logrus.Logger
	logFolder         = getString("log.folder")
	enableFileLogging = getBoolOrDefault("log.file.enable", true)
	filename          = getStringOrDefault("log.file.name", fmt.Sprintf("%s%s", "application", logExtension))
)

func init() {
	if !strings.HasSuffix(filename, logExtension) {
		filename = fmt.Sprintf("%s%s", filename, logExtension)
	}
	log = logrus.New()
	//log.SetFormatter(&logrus.TextFormatter{})
	log.SetFormatter(&Formatter{})
	log.SetLevel(getLogLevel())
	var mw io.Writer = os.Stdout
	if enableFileLogging {
		err := os.MkdirAll(logFolder, os.ModePerm)
		if err != nil {
			logging.Fatalf("Error creating log folder: %v", err)
		}

		mw = io.MultiWriter(os.Stdout, &lumberjack.Logger{
			Filename:   fmt.Sprintf("%s/%s", logFolder, filename),
			MaxSize:    configuration.GetIntOrDefault("log.file.maxSize", 10),     // Max size in MB
			MaxBackups: configuration.GetIntOrDefault("log.file.maxBackups", 3),   // Max number of old log files to keep
			MaxAge:     configuration.GetIntOrDefault("log.file.maxAge", 28),      // Max age in days to keep a log file
			Compress:   configuration.GetBoolOrDefault("log.file.compress", true), // Compress old log files
		})
	}
	log.SetOutput(mw)
}

// Info ...
func Info(v ...interface{}) {
	fields := getFields(runtime.Caller(1))
	log.WithFields(fields).Infoln(v...)
}

// Warn ...
func Warn(v ...interface{}) {
	fields := getFields(runtime.Caller(1))
	log.WithFields(fields).Warnln(v...)
}

// Debug ...
func Debug(v ...interface{}) {
	fields := getFields(runtime.Caller(1))
	log.WithFields(fields).Debugln(v...)
}

// Trace ...
func Trace(v ...interface{}) {
	fields := getFields(runtime.Caller(1))
	log.WithFields(fields).Traceln(v...)
}

// Error ...
func Error(v ...interface{}) {
	fields := getFields(runtime.Caller(1))
	log.WithFields(fields).Errorln(v...)
}

// Fatal ...
func Fatal(v ...interface{}) {
	fields := getFields(runtime.Caller(1))
	log.WithFields(fields).Fatalln(v...)
}

// Infof ...
func Infof(format string, v ...interface{}) {
	fields := getFields(runtime.Caller(1))
	log.WithFields(fields).Infof(format, v...)
}

// Warnf ...
func Warnf(format string, v ...interface{}) {
	fields := getFields(runtime.Caller(1))
	log.WithFields(fields).Warnf(format, v...)
}

// Debugf ...
func Debugf(format string, v ...interface{}) {
	fields := getFields(runtime.Caller(1))
	log.WithFields(fields).Debugf(format, v...)
}

// Tracef ...
func Tracef(format string, v ...interface{}) {
	fields := getFields(runtime.Caller(1))
	log.WithFields(fields).Tracef(format, v...)
}

// Errorf ...
func Errorf(format string, v ...interface{}) {
	fields := getFields(runtime.Caller(1))
	log.WithFields(fields).Errorf(format, v...)
}

// Fatalf ...
func Fatalf(format string, v ...interface{}) {
	fields := getFields(runtime.Caller(1))
	log.WithFields(fields).Fatalf(format, v...)
}

var (

	// ConfigError ...
	ConfigError = "%v type=config.error"

	// HTTPError ...
	HTTPError = "%v type=http.error"

	// HTTPWarn ...
	HTTPWarn = "%v type=http.warn"

	// HTTPInfo ...
	HTTPInfo = "%v type=http.info"
)

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

func getLogLevel() logrus.Level {
	lvl := strings.ToLower(getStringOrDefault("log.level", "info"))
	switch lvl {
	case "info":
		return logrus.InfoLevel
	case "debug":
		return logrus.DebugLevel
	case "trace":
		return logrus.TraceLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	default:
		return logrus.InfoLevel

	}
}
func getFields(pc uintptr, file string, line int, ok bool) logrus.Fields {
	var fields = logrus.Fields{"file": "unknown"}
	if ok {
		fields = logrus.Fields{
			"file": filepath.Base(file),
			"line": line,
		}
	}
	return fields
}
