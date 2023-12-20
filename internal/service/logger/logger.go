package logger

import (
	"context"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

type Params struct {
	FilePath   string
	Level      string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

type Logger struct {
	*slog.Logger
}

const LevelFatal = slog.Level(12)

func GetLogger(params Params) (Logger, error) {
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     getLevel(params.Level),
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key != slog.TimeKey {
				return a
			}
			t := a.Value.Time()
			a.Value = slog.StringValue(t.Format(time.DateTime))
			return a
		},
	}

	if err := checkFilePermissions(params.FilePath); err != nil {
		return Logger{}, err
	}

	fileWriter := &lumberjack.Logger{
		Filename:   params.FilePath,
		MaxSize:    params.MaxSize,
		MaxBackups: params.MaxBackups,
		MaxAge:     params.MaxAge,
		Compress:   params.Compress,
	}

	mw := io.MultiWriter(os.Stdout, fileWriter)

	return Logger{slog.New(slog.NewTextHandler(mw, opts))}, nil
}

func checkFilePermissions(filePath string) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0644); err != nil {
		return err
	}
	if f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644); err != nil {
		return err
	} else if err := f.Close(); err != nil {
		return err
	}

	return nil
}

func (l Logger) Fatal(msg string, args ...any) {
	ctx := context.Background()
	l.Log(ctx, LevelFatal, msg, args)
	os.Exit(1)
}

func (l Logger) DebugOrError(err error, msg string, args ...any) {
	if err == nil {
		l.Debug(msg, args)
	} else {
		l.Error(msg, args)
	}
}

func getLevel(level string) slog.Level {
	switch level {
	case "fatal":
		return LevelFatal
	case "error":
		return slog.LevelError
	case "warn":
		return slog.LevelWarn
	case "info":
		return slog.LevelInfo
	case "debug":
		return slog.LevelDebug
	}

	return slog.LevelInfo
}
