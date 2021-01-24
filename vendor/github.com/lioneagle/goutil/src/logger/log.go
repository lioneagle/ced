package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"sync"
	"time"
)

const STACK_BUFFER_SIZE int = 8192
const TIMESTAMP_FMT string = "2006-01-02 15:04:05.000"

const (
	EMERGENCY = 0
	ALERT     = 1
	CRITICAL  = 2
	ERROR     = 3
	WARNING   = 4
	NOTICE    = 5
	INFO      = 6
	DEBUG     = 7
)

type Level int

var levelNames = []string{"EMERGENCY", "ALERT", "CRITICAL", "ERROR", "WARNING", "NOTICE", "INFO", "DEBUG"}

func (level Level) String() string {
	if level < 0 || int(level) >= len(levelNames) {
		return "???"
	}
	return levelNames[level]
}

type Logger struct {
	showShortFile      bool
	showPackage        bool
	showFuncName       bool
	shortFileNameDepth int
	out                io.Writer  // destination for output
	mu                 sync.Mutex // ensures atomic writes
	level              Level
	stackTraceLevel    Level
	modules            map[int]string
}

func NewLogger(out io.Writer) *Logger {
	logger := &Logger{out: out}
	logger.init()
	return logger
}

func (this *Logger) init() {
	this.level = WARNING
	this.stackTraceLevel = EMERGENCY
	this.showShortFile = true
	this.shortFileNameDepth = 1
	this.showPackage = true
	this.showFuncName = false
}

func (this *Logger) SetLevel(level Level) {
	this.level = level
}

func (this *Logger) SetStackTraceLevel(level Level) {
	this.stackTraceLevel = level
}

func (this *Logger) ShowShortFileName() {
	this.showShortFile = true
}

func (this *Logger) SetShortFileNameDepth(depth int) {
	this.shortFileNameDepth = depth
}

func (this *Logger) ShowFullFileName() {
	this.showShortFile = false
}

func (this *Logger) ShowFuncName() {
	this.showFuncName = true
}

func (this *Logger) HideFuncName() {
	this.showFuncName = false
}

func (this *Logger) ShowPackage() {
	this.showPackage = true
}

func (this *Logger) HidePackage() {
	this.showPackage = false
}

func (this *Logger) log(level Level, msg string, args ...interface{}) {
	if level > this.level {
		return
	}

	// Write all the data into a buffer.
	// Format is:
	// <timestamp> [level][<file>:<line>:<function>]: <message>
	now := time.Now()

	var buffer bytes.Buffer
	buffer.WriteString(now.Format(TIMESTAMP_FMT))
	buffer.WriteString(" ")
	buffer.WriteString(fmt.Sprintf("[%s]", level.String()))
	buffer.WriteString(this.fileInfo(4))
	buffer.WriteString(": ")
	buffer.WriteString(fmt.Sprintf(msg, args...))
	buffer.WriteString("\n")

	if level <= this.stackTraceLevel {
		buffer.WriteString("--- BEGIN stacktrace: ---\n")
		buffer.Write(stackTrace())
		buffer.WriteString("--- END stacktrace ---\n\n")
	}

	this.output(buffer.Bytes())
}

func (this *Logger) print(format string, args ...interface{}) {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf(format, args...))
	buffer.WriteString("\n")
	this.output(buffer.Bytes())
}

func (this *Logger) output(msg []byte) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.out.Write(msg)
}

func (this *Logger) fileInfo(depth int) string {
	stackInfo := "[???]"
	if pc, fileName, line, ok := runtime.Caller(depth); ok {

		if this.showShortFile {
			fileName = extractFileName(fileName, this.shortFileNameDepth)
		}

		if this.showFuncName {
			funcName := runtime.FuncForPC(pc).Name()
			if !this.showPackage {
				funcName = path.Base(funcName)
			}
			stackInfo = fmt.Sprintf("[%s:%d:%s]", fileName, line, funcName)
		} else {
			stackInfo = fmt.Sprintf("[%s:%d]", fileName, line)
		}

	}
	return stackInfo
}

func stackTrace() []byte {
	trace := make([]byte, STACK_BUFFER_SIZE)
	count := runtime.Stack(trace, true)
	return trace[:count]
}

func extractFileName(fileName string, shortFileNameDepth int) string {
	if shortFileNameDepth < 1 {
		shortFileNameDepth = 1
	}

	i := 0
	depth := 0
	for i = len(fileName) - 1; i >= 0; i-- {
		if fileName[i] == '/' {
			depth++
			if depth >= shortFileNameDepth {
				break
			}
		}
	}
	return fileName[i+1 : len(fileName)]
}

func (this *Logger) Emergency(msg string) {
	this.log(EMERGENCY, "%s", msg)
}

func (this *Logger) Emergencyf(format string, args ...interface{}) {
	this.log(EMERGENCY, format, args...)
}

func (this *Logger) Alert(msg string) {
	this.log(ALERT, "%s", msg)
}

func (this *Logger) Alertf(format string, args ...interface{}) {
	this.log(ALERT, format, args...)
}

func (this *Logger) Critical(msg string) {
	this.log(CRITICAL, "%s", msg)
}

func (this *Logger) Criticalf(format string, args ...interface{}) {
	this.log(CRITICAL, format, args...)
}

func (this *Logger) Error(msg string) {
	this.log(ERROR, "%s", msg)
}

func (this *Logger) Errorf(format string, args ...interface{}) {
	this.log(ERROR, format, args...)
}

func (this *Logger) Warning(msg string) {
	this.log(WARNING, "%s", msg)
}

func (this *Logger) Warningf(format string, args ...interface{}) {
	this.log(WARNING, format, args...)
}

func (this *Logger) Notice(msg string) {
	this.log(NOTICE, "%s", msg)
}

func (this *Logger) Noticef(format string, args ...interface{}) {
	this.log(NOTICE, format, args...)
}

func (this *Logger) Info(msg string) {
	this.log(INFO, "%s", msg)
}

func (this *Logger) Infof(format string, args ...interface{}) {
	this.log(INFO, format, args...)
}

func (this *Logger) Debug(msg string) {
	this.log(DEBUG, "%s", msg)
}

func (this *Logger) Debugf(format string, args ...interface{}) {
	this.log(DEBUG, format, args...)
}

func (this *Logger) Print(msg string) {
	this.print("%s", msg)
}

func (this *Logger) Printf(format string, args ...interface{}) {
	this.print(format, args...)
}

func (this *Logger) PrintStack() {
	this.output(stackTrace())
}

var defaultLogger *Logger = nil

func getDefaultLogger() *Logger {
	if defaultLogger != nil {
		return defaultLogger
	}
	defaultLogger = NewLogger(os.Stderr)
	return defaultLogger
}

func Emergency(msg string) {
	getDefaultLogger().Emergency(msg)
}

func Emergencyf(format string, args ...interface{}) {
	getDefaultLogger().Emergencyf(format, args...)
}

func Alert(msg string) {
	getDefaultLogger().Alert(msg)
}

func Alertf(format string, args ...interface{}) {
	getDefaultLogger().Alertf(format, args...)
}

func Critical(msg string) {
	getDefaultLogger().Critical(msg)
}

func Criticalf(format string, args ...interface{}) {
	getDefaultLogger().Criticalf(format, args...)
}

func Error(msg string) {
	getDefaultLogger().Error(msg)
}

func Errorf(format string, args ...interface{}) {
	getDefaultLogger().Errorf(format, args...)
}

func Warning(msg string) {
	getDefaultLogger().Warning(msg)
}

func Warningf(format string, args ...interface{}) {
	getDefaultLogger().Warningf(format, args...)
}

func Notice(msg string) {
	getDefaultLogger().Notice(msg)
}

func Noticef(msg string, args ...interface{}) {
	getDefaultLogger().Noticef(msg, args...)
}

func Info(msg string) {
	getDefaultLogger().Info(msg)
}

func Infof(format string, args ...interface{}) {
	getDefaultLogger().Infof(format, args...)
}

func Debug(msg string) {
	getDefaultLogger().Debug(msg)
}

func Debugf(format string, args ...interface{}) {
	getDefaultLogger().Debugf(format, args...)
}

func Print(msg string) {
	getDefaultLogger().Print(msg)
}

func Printf(format string, args ...interface{}) {
	getDefaultLogger().Printf(format, args...)
}

func PrintStack() {
	getDefaultLogger().PrintStack()
}

func SetLevel(level Level) {
	getDefaultLogger().SetLevel(level)
}

func SetStackTraceLevel(level Level) {
	getDefaultLogger().SetStackTraceLevel(level)
}

func ShowShortFileName() {
	getDefaultLogger().ShowShortFileName()
}

func SetShortFileNameDepth(depth int) {
	getDefaultLogger().SetShortFileNameDepth(depth)
}

func ShowFullFileName() {
	getDefaultLogger().ShowFullFileName()
}

func ShowFuncName() {
	getDefaultLogger().ShowFuncName()
}

func HideFuncName() {
	getDefaultLogger().HideFuncName()
}

func ShowPackage() {
	getDefaultLogger().ShowPackage()
}

func HidePackage() {
	getDefaultLogger().HidePackage()
}
