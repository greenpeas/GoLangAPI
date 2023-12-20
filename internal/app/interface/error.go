package app_interface

type Error interface {
	Error() string
	Unwrap() error
	GetBody() any
	GetHttpCode() int
}
