package customtags

import (
	"reflect"
)

// __normalize returns correct reflect values according to its kind.
func (c *Impl) __normalize(k reflect.Kind, v reflect.Value) reflect.Value {
	switch k {
	case reflect.Struct:
		if v.Kind() == reflect.Ptr {
			return v
		} else {
			return v.Elem()
		}
	case reflect.Map, reflect.Slice, reflect.Array, reflect.Ptr:
		return v
	default:
		return v.Elem()
	}
}

// __parse parse data recursively, modify fields data if custom tag labels presented.
func (c *Impl) __parse(field string, v reflect.Value, tag string) reflect.Value {
	var (
		orig  = v
		cpVal reflect.Value
	)

	// if input is pointer, dereference it properly.
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return reflect.New(v.Type()).Elem()
		}

		orig = v.Elem()
	}

	// initial tag = "", but since initial value v could be struct or *struct only,
	// first iteration always get to Struct case.
	switch orig.Kind() {
	case reflect.Struct:
		// copy type of original entity.
		t := orig.Type()
		// get new ptr of created entity with original entity's type.
		cpVal = reflect.New(t)

		// iterate over struct fields.
		for i := 0; i < t.NumField(); i++ {
			// get copied struct type's field i.
			f := t.Field(i)
			// get original struct type's field i.
			fVal := orig.Field(i)

			// if we won't be able to get Interface of field, skip it to prevent panic.
			if !fVal.CanInterface() {
				continue
			}

			// get custom tag label.
			label := f.Tag.Get(c.tag)
			if label == "" {
				label = tag
			}

			// if field type's kind is ptr or string, cast modified data to interface{} and then to initial type,
			// only after type casting set modified value to copied struct field.
			if fVal.Type().Kind() == reflect.Ptr {
				fVal = fVal.Elem()
			}

			if fVal.Kind() == reflect.String {
				s := c.__parse(f.Name, fVal, label).Interface().(string)
				// if copied Struct is a Pointer, set the string properly.
				if cpVal.Kind() == reflect.Ptr {
					cpVal.Elem().Field(i).SetString(s)
				} else {
					cpVal.Field(i).SetString(s)
				}
			} else if fVal.Kind() == reflect.Struct {
				modVal := c.__parse(f.Name, fVal, label)
				if modVal.Kind() == reflect.Ptr {
					if cpVal.Elem().Field(i).Kind() == reflect.Ptr {
						cpVal.Elem().Field(i).Set(modVal)
					} else {
						cpVal.Elem().Field(i).Set(modVal.Elem())
					}
				} else {
					cpVal.Field(i).Set(modVal)
				}
			} else {
				if v.Kind() == reflect.Ptr {
					cpVal.Field(i).Set(c.__parse(f.Name, fVal, label))
				} else {
					cpVal.Elem().Field(i).Set(c.__parse(f.Name, fVal, label))
				}
			}
		}
	case reflect.Slice, reflect.Array:
		cpVal = reflect.MakeSlice(orig.Type(), orig.Len(), orig.Cap())

		// slice/array values could not have tags.
		for i := 0; i < orig.Len(); i++ {
			cpVal.Index(i).Set(c.__parse(field, orig.Index(i), ""))
		}
	case reflect.Map:
		cpVal = reflect.MakeMap(orig.Type())
		keys := orig.MapKeys()

		// map values could not have tags.
		for i := 0; i < orig.Len(); i++ {
			cpVal.SetMapIndex(keys[i], c.__parse(keys[i].String(), orig.MapIndex(keys[i]), ""))
		}
	default:
		if reflect.ValueOf(orig).IsZero() {
			cpVal = orig
			break
		}

		cpVal = reflect.New(orig.Type())

		if v.Kind() == reflect.Ptr && !v.IsNil() {
			v = v.Elem()
		}

		if modVal, ok := c.__handle(v.Interface(), tag); ok {
			cpVal.Elem().Set(reflect.ValueOf(modVal))
		} else {
			cpVal.Elem().Set(orig)
		}
	}

	return c.__normalize(v.Kind(), cpVal)
}

// __handle returns field's value parsed with Handler according to tag label.
// If it is empty, returns initial value.
func (c *Impl) __handle(input any, tag string) (any, bool) {
	if container, ok := customTags.getHandler(tag); ok {
		if reflect.ValueOf(input).Kind() == container.Kind {
			return container.Handler(input), true
		}
	}

	return input, false
}
