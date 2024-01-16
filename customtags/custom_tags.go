package customtags

import (
	"reflect"
	"sync"
)

// Handler is a func(any) any.
type Handler func(any) any

type Container struct {
	Handler Handler
	Kind    reflect.Kind
}

// Mapper is a map[label]handler_func.
type Mapper map[string]Container

type Impl struct {
	mx      sync.Mutex
	tag     string
	mappers Mapper
}

func newCustomTagsImpl(tag string) *Impl {
	return &Impl{tag: tag}
}

type CustomTagger interface {
	bind(key string, fn Handler, kind reflect.Kind)
	getHandler(tag string) (fn Container, ok bool)

	Proceed(input any) (output any)
}

var customTags CustomTagger

// InitCustomTags inits global customTags.
// Example: customtags.InitCustomTags(key); customtags.customTags.Proceed(model)
func InitCustomTags(key string) {
	customTags = newCustomTagsImpl(key)
}

func CT() CustomTagger {
	return customTags
}

// NewCustomTags returns CustomTagger instance for Instant or Dependency Injection usage.
// Example: sm := customtags.NewCustomTags(key, true); sm.Proceed(model)
// If showStackTrace is true, Proceed method prints stack trace of recovered panic.
func NewCustomTags(key string) *Impl {
	sm := newCustomTagsImpl(key)

	customTags = sm

	return sm
}

// bind connects Handler with tag label.
func (c *Impl) bind(label string, fn Handler, kind reflect.Kind) {
	if c.mappers == nil {
		c.mappers = make(map[string]Container, 1)
	}

	c.mx.Lock()
	c.mappers[label] = Container{
		Handler: fn,
		Kind:    kind,
	}
	c.mx.Unlock()
}

// getHandler returns Handler or nil
func (c *Impl) getHandler(tag string) (Container, bool) {
	handler, ok := c.mappers[tag]

	return handler, ok
}

// Bind connects custom handler of any type with tag label.
// Example: customtags.Bind("custom_label", func(string) string {...})
func Bind[T any](label string, fn func(T) T) {
	if customTags == nil {
		return
	}

	var t T
	afn := func(in any) any {
		return fn(in.(T))
	}

	customTags.bind(label, afn, reflect.TypeOf(t).Kind())
}

// Proceed replace tagged fields of the struct with modified by the Handler data.
func (c *Impl) Proceed(input any) any {
	if input == nil ||
		(reflect.ValueOf(input).Kind() != reflect.Struct &&
			reflect.ValueOf(input).Kind() != reflect.Ptr) {
		return nil
	}

	ret := c.__parse("", reflect.ValueOf(input), "")

	if !ret.IsValid() {
		return input
	}

	return ret.Interface()
}
