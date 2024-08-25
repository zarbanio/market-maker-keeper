package logger

import (
	"context"
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/zarbanio/market-maker-keeper/store"
)

var Logger zerolog.Logger

type CustomWriter struct {
	ctx    context.Context
	writer io.Writer
	istore store.IStore
}

func NewCustomWriter(ctx context.Context, istore store.IStore, writer io.Writer) *CustomWriter {
	return &CustomWriter{
		ctx:    ctx,
		writer: writer,
		istore: istore,
	}
}

func (cw *CustomWriter) Write(p []byte) (n int, err error) {
	_, err = cw.istore.CreateLog(cw.ctx, p)
	if err != nil {
		return 0, err
	}

	return cw.writer.Write(p)
}

func InitLogger(ctx context.Context, s store.IStore, level string) error {
	customWriter := NewCustomWriter(ctx, s, os.Stdout)
	lvl := ParseLevel(level)

	Logger = zerolog.New(customWriter).With().Logger().Output(customWriter).Level(lvl)
	return nil
}

func ParseLevel(level string) zerolog.Level {
	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}
