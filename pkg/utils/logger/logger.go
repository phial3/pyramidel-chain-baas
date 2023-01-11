// Copyright (c) 2022 s1ren
// hxx258456/pyramidel-chain-baas is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 			http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hxx258456/pyramidel-chain-baas/pkg/constants"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	defaultLogger *zap.Logger
)

type BLogger struct {
	sugard *zap.SugaredLogger
}

// init default logger with only console output info above
func init() {
	zc := zapcore.NewTee(newConsoleCore(zap.InfoLevel))
	defaultLogger = zap.New(zc)
}

func NewLogger(name string) BLogger {
	return BLogger{
		sugard: defaultLogger.Named(name).Sugar(),
	}
}

// CfgConsoleLogger config for console logs
// cfg donot support concurrent calls (as any package should init cfg at startup once)
func CfgConsoleLogger(debugMode bool, showPath bool) {
	level, zos := genConfigs(debugMode, showPath)

	zc := zapcore.NewTee(newConsoleCore(level))
	defaultLogger = zap.New(zc, zos...)
}

// TODO: export more file configs
// CfgConsoleAndFileLogger config for both console and file logs
// cfg donot support concurrent calls (as any package should init cfg at startup once)
func CfgConsoleAndFileLogger(debugMode bool, name string, showPath bool) {
	level, zos := genConfigs(debugMode, showPath)

	zc := zapcore.NewTee(newConsoleCore(level), newFileCore(name, level))

	defaultLogger = zap.New(zc, zos...)
}

func genConfigs(debugMode bool, showPath bool) (zapcore.LevelEnabler, []zap.Option) {
	level := zapcore.InfoLevel
	if debugMode {
		level = zapcore.DebugLevel
	}

	zos := []zap.Option{
		// zap.AddStacktrace(zapcore.WarnLevel),
	}
	if showPath {
		// skip self wrapper
		zos = append(zos, zap.AddCaller(), zap.AddCallerSkip(1))
	}

	return level, zos
}

func newConsoleCore(le zapcore.LevelEnabler) zapcore.Core {
	consoleLogger := zapcore.Lock(os.Stdout)

	zec := zap.NewProductionEncoderConfig()
	zec.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zec.EncodeTime = zapcore.ISO8601TimeEncoder
	zec.EncodeTime = shortTimeEncoder
	// zec.EncodeTime = zapcore.ISO8601TimeEncoder
	zec.ConsoleSeparator = " "

	consoleEncoder := zapcore.NewConsoleEncoder(zec)

	return zapcore.NewCore(consoleEncoder, consoleLogger, le)
}

func newFileCore(filename string, le zapcore.LevelEnabler) zapcore.Core {
	//TODO: export more rotate configs
	fileLogger := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filename,
		MaxSize:    1, // megabytes per file
		Compress:   true,
		MaxAge:     14,
		MaxBackups: 20,
		LocalTime:  true,
	})
	zec := zap.NewProductionEncoderConfig()
	zec.EncodeTime = zapcore.ISO8601TimeEncoder

	fileEncoder := zapcore.NewJSONEncoder(zec)
	return zapcore.NewCore(fileEncoder, fileLogger, le)
}

func shortTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(constants.ShortTimeLayout))
}

// IsDebugMode check DebugLevel enabled
func IsDebugMode() bool {
	return defaultLogger.Core().Enabled(zapcore.DebugLevel)
}

// Fatal logs a message at emergency level and exit.
func Fatal(f interface{}, v ...interface{}) {
	defaultLogger.Sugar().Fatalf(formatLog(zapcore.FatalLevel, f, v...))
}

func (b *BLogger) Fatal(f interface{}, v ...interface{}) {
	b.sugard.Fatalf(formatLog(zapcore.FatalLevel, f, v...))
}

// Panic logs a message at emergency level and exit.
func Panic(f interface{}, v ...interface{}) {
	defaultLogger.Sugar().Panicf(formatLog(zapcore.PanicLevel, f, v...))
}

func (b *BLogger) Panic(f interface{}, v ...interface{}) {
	b.sugard.Panicf(formatLog(zapcore.PanicLevel, f, v...))
}

// Error logs a message at error level.
func Error(f interface{}, v ...interface{}) {
	defaultLogger.Sugar().Errorf(formatLog(zapcore.ErrorLevel, f, v...))
}

// Error logs a message at error level.
func (b *BLogger) Error(f interface{}, v ...interface{}) {
	b.sugard.Errorf(formatLog(zapcore.ErrorLevel, f, v...))
}

// Warn logs a message at warning level.
func Warn(f interface{}, v ...interface{}) {
	defaultLogger.Sugar().Warnf(formatLog(zapcore.WarnLevel, f, v...))
}

func (b *BLogger) Warn(f interface{}, v ...interface{}) {
	b.sugard.Warnf(formatLog(zapcore.WarnLevel, f, v...))
}

// Info logs a message at info level.
func Info(f interface{}, v ...interface{}) {
	defaultLogger.Sugar().Infof(formatLog(zapcore.InfoLevel, f, v...))
}

func (b *BLogger) Info(f interface{}, v ...interface{}) {
	b.sugard.Infof(formatLog(zapcore.InfoLevel, f, v...))
}

// Debug logs a message at debug level.
func Debug(f interface{}, v ...interface{}) {
	defaultLogger.Sugar().Debugf(formatLog(zapcore.DebugLevel, f, v...))
}

func (b *BLogger) Debug(f interface{}, v ...interface{}) {
	b.sugard.Debugf(formatLog(zapcore.DebugLevel, f, v...))
}

func formatLog(l zapcore.Level, f interface{}, v ...interface{}) string {
	var msg string
	switch f := f.(type) {
	case string:
		msg = f
		if len(v) == 0 {
			return appendColor(l, msg)
		}
		if strings.Contains(msg, "%") && !strings.Contains(msg, "%%") {
			//format string
		} else {
			//do not contain format char
			msg += strings.Repeat(" %v", len(v))
		}
	default:
		msg = fmt.Sprint(f)
		if len(v) == 0 {
			return appendColor(l, msg)
		}
		msg += strings.Repeat(" %v", len(v))
	}
	return appendColor(l, fmt.Sprintf(msg, v...))
}

func appendColor(l zapcore.Level, s string) string {
	// default all red
	return s
	// println(1)
	// c := uint8(31)
	// switch l {
	// case zapcore.DebugLevel:
	// 	c = uint8(35) // Magenta
	// case zapcore.InfoLevel:
	// 	c = uint8(34) // Blue
	// case zapcore.WarnLevel:
	// 	c = uint8(33) // Yellow
	// }
	// return fmt.Sprintf("\x1b[%dm%s\x1b[0m", c, s)
}
