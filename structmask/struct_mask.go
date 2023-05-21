package structmask

import (
	"reflect"
)

// Handler is a func(any) any.
type Handler func(any) any

// Mapper is a map[label]handler_func.
type Mapper map[string]Handler

type SM struct {
	tag     string
	mappers Mapper
}

func newSM(tag string) *SM {
	return &SM{tag: tag}
}

type StructMask interface {
	bind(key string, fn Handler)
	getHandler(tag string) (fn Handler, ok bool)

	Modify(input any) (output any)
}

var StructMasker StructMask

// InitStructMask inits global StructMasker.
// Example: structmask.InitStructMask(cfg); structmask.StructMasker.Modify(model)
func InitStructMask(key string) {
	StructMasker = newSM(key)
}

// NewStructMask returns StructMask instance for Instant or Dependency Injection usage.
// Example: sm := structmask.NewStructMask(cfg); sm.Modify(model)
func NewStructMask(key string) *SM {
	sm := newSM(key)

	StructMasker = sm

	return sm
}

// bind connects Handler with tag label.
func (sm *SM) bind(label string, fn Handler) {
	if sm.mappers == nil {
		sm.mappers = make(map[string]Handler, 1)
	}

	sm.mappers[label] = fn
}

// getHandler returns Handler or nil
func (sm *SM) getHandler(tag string) (Handler, bool) {
	handler, ok := sm.mappers[tag]

	return handler, ok
}

// Bind connects custom handler of any type with tag label.
// Example: structmask.Bind("custom_label", func(string) string {...})
func Bind[T any](label string, fn func(T) T) {
	if StructMasker == nil {
		return
	}

	afn := func(in any) any {
		return fn(in.(T))
	}

	StructMasker.bind(label, afn)
}

// Modify replace tagged fields of the struct with modified by the Handler data.
func (sm *SM) Modify(input any) any {
	if input == nil ||
		(reflect.ValueOf(input).Kind() != reflect.Struct &&
			reflect.ValueOf(input).Kind() != reflect.Ptr) {
		return nil
	}

	return sm.__parse("", reflect.ValueOf(input), "").Interface()
}
