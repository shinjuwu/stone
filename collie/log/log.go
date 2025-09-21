package log

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	infoLogger *logrus.Logger
	infoFile   *rotatelogs.RotateLogs
	errLogger  *logrus.Logger
	errFile    *rotatelogs.RotateLogs
}

func New(pathname string, Level int32) (*Logger, error) {
	// new
	logger := new(Logger)
	logger.infoLogger, logger.infoFile = CreatLog(pathname, "Info")
	logger.infoLogger.SetLevel(logrus.Level(Level))
	logger.errLogger, logger.errFile = CreatLog(pathname, "Err")
	//Setting
	logger.infoLogger.SetOutput(io.MultiWriter(os.Stdout, logger.infoFile))
	logger.errLogger.SetOutput(io.MultiWriter(os.Stdout, logger.errFile))

	return logger, nil
}
func CreatLog(pathname string, typename string) (*logrus.Logger, *rotatelogs.RotateLogs) {
	var baseLogger *logrus.Logger
	var baseFile *rotatelogs.RotateLogs
	if pathname != "" {
		path := path.Join(filepath.Dir(pathname), "GameHub"+typename+".log")
		fileWriter, _ := rotatelogs.New(
			path+".%Y%m%d",
			rotatelogs.WithLinkName(path),
			rotatelogs.WithMaxAge(time.Duration(336)*time.Hour),
			rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
		)
		baseLogger = logrus.New()
		baseFile = fileWriter

	} else {
		baseLogger = logrus.New()
	}
	return baseLogger, baseFile
}

// It's dangerous to call the method on logging
func (logger *Logger) Close() {
	if logger.infoFile != nil {
		logger.infoFile.Close()
	}
	if logger.errFile != nil {
		logger.errFile.Close()
	}

	logger.infoLogger = nil
	logger.infoFile = nil
	logger.errLogger = nil
	logger.errFile = nil
}

func (logger *Logger) Trace(format string, a ...interface{}) {
	logger.infoLogger.Tracef(format, a...)
}

func (logger *Logger) Debug(format string, a ...interface{}) {
	logger.infoLogger.Debugf(format, a...)
}

func (logger *Logger) Info(format string, a ...interface{}) {
	logger.infoLogger.Infof(format, a...)
}

func (logger *Logger) Release(format string, a ...interface{}) {
	logger.infoLogger.Infof(format, a...)
}

func (logger *Logger) Error(format string, a ...interface{}) {
	logger.errLogger.Errorf(format, a...)
}

func (logger *Logger) Fatal(format string, a ...interface{}) {
	logger.errLogger.Fatalf(format, a...)
}

var gLogger *Logger

// It's dangerous to call the method on logging
func Export(logger *Logger) {
	if logger != nil {
		gLogger = logger
	}
}

func Trace(format string, a ...interface{}) {
	gLogger.infoLogger.Tracef(format, a...)
}

func Debug(format string, a ...interface{}) {
	gLogger.infoLogger.Debugf(format, a...)
}

func Release(format string, a ...interface{}) {
	gLogger.infoLogger.Infof(format, a...)
}

func Info(format string, a ...interface{}) {
	gLogger.infoLogger.Infof(format, a...)
}

func Error(format string, a ...interface{}) {
	gLogger.errLogger.Errorf(format, a...)
}

func Fatal(format string, a ...interface{}) {
	gLogger.errLogger.Fatalf(format, a...)
}

func Close() {
	gLogger.Close()
}
