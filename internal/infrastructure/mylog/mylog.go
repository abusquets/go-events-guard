// logger.go
package mylog

type Logger interface {
	Info(msg string)
	Error(msg string)
	ErrorWithErr(msg string, err error)
	Debug(msg string)
	Warn(msg string)
	With(key string, value interface{}) Logger
}
