package app_interface

type Logger interface {
	Fatal(msg string, args ...any)
	Error(msg string, args ...any)
	Info(msg string, args ...any)
	Debug(msg string, args ...any)
	DebugOrError(err error, msg string, args ...any)
}
