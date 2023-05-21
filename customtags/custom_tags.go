package customtags

import (
	"reflect"
)

// Handler is a func(any) any.
type Handler func(any) any

// Mapper is a map[label]handler_func.
type Mapper map[string]Handler

type CT struct {
	tag     string
	mappers Mapper
}

func newSM(tag string) *CT {
	return &CT{tag: tag}
}

type CustomTags interface {
	bind(key string, fn Handler)
	getHandler(tag string) (fn Handler, ok bool)

	Modify(input any) (output any)
}

var CustomTagger CustomTags

// InitCustomTags inits global CustomTagger.
// Example: customtags.InitCustomTags(cfg); customtags.CustomTagger.Modify(model)
func InitCustomTags(key string) {
	CustomTagger = newSM(key)
}

// NewCustomTags returns CustomTags instance for Instant or Dependency Injection usage.
// Example: sm := customtags.NewCustomTags(cfg); sm.Modify(model)
func NewCustomTags(key string) *CT {
	sm := newSM(key)

	CustomTagger = sm

	return sm
}

// bind connects Handler with tag label.
func (c *CT) bind(label string, fn Handler) {
	if c.mappers == nil {
		c.mappers = make(map[string]Handler, 1)
	}

	c.mappers[label] = fn
}

// getHandler returns Handler or nil
func (c *CT) getHandler(tag string) (Handler, bool) {
	handler, ok := c.mappers[tag]

	return handler, ok
}

// Bind connects custom handler of any type with tag label.
// Example: customtags.Bind("custom_label", func(string) string {...})
func Bind[T any](label string, fn func(T) T) {
	if CustomTagger == nil {
		return
	}

	afn := func(in any) any {
		return fn(in.(T))
	}

	CustomTagger.bind(label, afn)
}

// Modify replace tagged fields of the struct with modified by the Handler data.
func (c *CT) Modify(input any) any {
	if input == nil ||
		(reflect.ValueOf(input).Kind() != reflect.Struct &&
			reflect.ValueOf(input).Kind() != reflect.Ptr) {
		return nil
	}

	return c.__parse("", reflect.ValueOf(input), "").Interface()
}
