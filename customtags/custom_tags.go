package customtags

import (
	"reflect"
)

// Handler is a func(any) any.
type Handler func(any) any

// Mapper is a map[label]handler_func.
type Mapper map[string]Handler

type CustomTagsImpl struct {
	tag     string
	mappers Mapper
}

func newCustomTagsImpl(tag string) *CustomTagsImpl {
	return &CustomTagsImpl{tag: tag}
}

type CustomTagger interface {
	bind(key string, fn Handler)
	getHandler(tag string) (fn Handler, ok bool)

	Proceed(input any) (output any)
}

var CustomTags CustomTagger

// InitCustomTags inits global CustomTags.
// Example: customtags.InitCustomTags(cfg); customtags.CustomTags.Proceed(model)
func InitCustomTags(key string) {
	CustomTags = newCustomTagsImpl(key)
}

// NewCustomTags returns CustomTagger instance for Instant or Dependency Injection usage.
// Example: sm := customtags.NewCustomTags(cfg); sm.Proceed(model)
func NewCustomTags(key string) *CustomTagsImpl {
	sm := newCustomTagsImpl(key)

	CustomTags = sm

	return sm
}

// bind connects Handler with tag label.
func (c *CustomTagsImpl) bind(label string, fn Handler) {
	if c.mappers == nil {
		c.mappers = make(map[string]Handler, 1)
	}

	c.mappers[label] = fn
}

// getHandler returns Handler or nil
func (c *CustomTagsImpl) getHandler(tag string) (Handler, bool) {
	handler, ok := c.mappers[tag]

	return handler, ok
}

// Bind connects custom handler of any type with tag label.
// Example: customtags.Bind("custom_label", func(string) string {...})
func Bind[T any](label string, fn func(T) T) {
	if CustomTags == nil {
		return
	}

	afn := func(in any) any {
		return fn(in.(T))
	}

	CustomTags.bind(label, afn)
}

// Proceed replace tagged fields of the struct with modified by the Handler data.
func (c *CustomTagsImpl) Proceed(input any) any {
	if input == nil ||
		(reflect.ValueOf(input).Kind() != reflect.Struct &&
			reflect.ValueOf(input).Kind() != reflect.Ptr) {
		return nil
	}

	return c.__parse("", reflect.ValueOf(input), "").Interface()
}
