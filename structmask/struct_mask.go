package structmask

import "reflect"

// Handler is a func(string)string.
type Handler func(string) string

// Mapper is a map[label]handler_func.
type Mapper map[string]Handler

type SM struct {
	cfg *Config
}

func newSM(cfg *Config) *SM {
	return &SM{cfg: cfg}
}

type Config struct {
	TagName string
	Mappers Mapper
}

type StructMask interface {
	Proceed(input any) (output any)
	getHandler(tag string) (fn Handler)
}

var StructMasker StructMask

// InitStructMask inits global StructMasker.
// Example: structmask.InitStructMask(cfg); structmask.StructMasker.Proceed(model)
func InitStructMask(cfg *Config) {
	StructMasker = newSM(cfg)
}

// NewStructMask returns StructMask instance for Instant or Dependency Injection usage.
// Example: sm := structmask.NewStructMask(cfg); sm.Proceed(model)
func NewStructMask(cfg *Config) *SM {
	return newSM(cfg)
}

// getHandler returns Handler or nil
func (sm *SM) getHandler(tag string) Handler {
	return sm.cfg.Mappers[tag]
}

// Proceed replace tagged fields of struct with correct masks
// Tag example: `mask:"cvv"`.
func (sm *SM) Proceed(input any) any {
	if input == nil ||
		(reflect.ValueOf(input).Kind() != reflect.Struct &&
			reflect.ValueOf(input).Kind() != reflect.Ptr) {
		return nil
	}

	return sm.__parse("", reflect.ValueOf(input), "").Interface()
}
