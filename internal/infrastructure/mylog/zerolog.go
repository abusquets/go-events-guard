package mylog

import (
	"os"
	"runtime"
	"strconv"
	"sync"

	"github.com/rs/zerolog"
)

var (
	singletonInstance *ZeroLogger
	once              sync.Once
)

// ZeroLogger és una implementació de la interfície Logger utilitzant zerolog.
type ZeroLogger struct {
	logger zerolog.Logger
}

// NewZeroLogger crea una nova instància de ZeroLogger amb configuració bàsica.
func newZeroLogger() *ZeroLogger {
	// Utilitza ColorWriter per a la sortida
	colorWriter := NewColorWriter(os.Stdout)
	zerologger := zerolog.New(colorWriter).With().Timestamp().Logger()
	return &ZeroLogger{logger: zerologger}
}

// GetLogger retorna una instància Singleton de ZeroLogger segura per a threads.
func GetLogger() *ZeroLogger {
	once.Do(func() {
		singletonInstance = newZeroLogger()
	})
	return singletonInstance
}

// Mètodes de la interfície Logger implementats per ZeroLogger
func (z *ZeroLogger) Info(msg string) {
	z.logger.Info().Str("caller", getCallerInfo()).Msg(msg)
}

func (z *ZeroLogger) Error(msg string) {
	z.logger.Error().Str("caller", getCallerInfo()).Msg(msg)
}

func (z *ZeroLogger) ErrorWithErr(msg string, err error) {
	z.logger.Error().Err(err).Str("caller", getCallerInfo()).Msg(msg)
}
func (z *ZeroLogger) Warn(msg string) {
	z.logger.Warn().Str("caller", getCallerInfo()).Msg(msg)
}
func (z *ZeroLogger) Debug(msg string) {
	z.logger.Debug().Str("caller", getCallerInfo()).Msg(msg)
}

func (z *ZeroLogger) With(key string, value interface{}) Logger {
	return &ZeroLogger{logger: z.logger.With().Interface(key, value).Logger()}
}

func getCallerInfo() string {
	// Obtenir el stack trace
	pc, file, line, ok := runtime.Caller(2) // 2 per saltar aquesta funció i la de log
	if !ok {
		return "unknown"
	}

	// Obtenir el nom de la funció
	funcName := runtime.FuncForPC(pc).Name()
	return funcName + " [" + file + ":" + strconv.Itoa(line) + "]"
}
