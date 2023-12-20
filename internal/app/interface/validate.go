package app_interface

type Validator interface {
	Struct(interface{}) map[string]string
}
