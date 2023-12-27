package customtags

import (
	"reflect"
	"sync"
)

// Handler is a func(any) any.
type Handler func(any) any

// Mapper is a map[label]handler_func.
type Mapper map[string]Handler

type Impl struct {
	mx             sync.Mutex
	tag            string
	mappers        Mapper
	showStackTrace bool // print stack trace on recovered panic
}

func newCustomTagsImpl(tag string, showStackTrace bool) *Impl {
	return &Impl{tag: tag, showStackTrace: showStackTrace}
}

type CustomTagger interface {
	bind(key string, fn Handler)
	getHandler(tag string) (fn Handler, ok bool)

	Proceed(input any) (output any)
}

var customTags CustomTagger

// InitCustomTags inits global customTags.
// Example: customtags.InitCustomTags(key); customtags.customTags.Proceed(model)
func InitCustomTags(key string, showStackTrace ...bool) {
	_showStackTrace := false
	if len(showStackTrace) > 0 {
		_showStackTrace = showStackTrace[0]
	}

	customTags = newCustomTagsImpl(key, _showStackTrace)
}

func CT() CustomTagger {
	return customTags
}

// NewCustomTags returns CustomTagger instance for Instant or Dependency Injection usage.
// Example: sm := customtags.NewCustomTags(key, true); sm.Proceed(model)
// If showStackTrace is true, Proceed method prints stack trace of recovered panic.
func NewCustomTags(key string, showStackTrace ...bool) *Impl {
	_showStackTrace := false
	if len(showStackTrace) > 0 {
		_showStackTrace = showStackTrace[0]
	}

	sm := newCustomTagsImpl(key, _showStackTrace)

	customTags = sm

	return sm
}

// bind connects Handler with tag label.
func (c *Impl) bind(label string, fn Handler) {
	if c.mappers == nil {
		c.mappers = make(map[string]Handler, 1)
	}

	c.mx.Lock()
	c.mappers[label] = fn
	c.mx.Unlock()
}

// getHandler returns Handler or nil
func (c *Impl) getHandler(tag string) (Handler, bool) {
	handler, ok := c.mappers[tag]

	return handler, ok
}

// Bind connects custom handler of any type with tag label.
// Example: customtags.Bind("custom_label", func(string) string {...})
func Bind[T any](label string, fn func(T) T) {
	if customTags == nil {
		return
	}

	afn := func(in any) any {
		return fn(in.(T))
	}

	customTags.bind(label, afn)
}

// Proceed replace tagged fields of the struct with modified by the Handler data.
func (c *Impl) Proceed(input any) any {
	if input == nil ||
		(reflect.ValueOf(input).Kind() != reflect.Struct &&
			reflect.ValueOf(input).Kind() != reflect.Ptr) {
		return nil
	}

	ret := c.__tryParse("", reflect.ValueOf(input), "")

	if !ret.IsValid() {
		return input
	}

	return ret.Interface()
}

func (c *Impl) MustProceed(input any) any {
	if input == nil ||
		(reflect.ValueOf(input).Kind() != reflect.Struct &&
			reflect.ValueOf(input).Kind() != reflect.Ptr) {
		return nil
	}

	return c.__parse("", reflect.ValueOf(input), "").Interface()
}
