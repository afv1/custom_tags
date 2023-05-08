package structmask

import (
	"reflect"
)

// __parse parse data recursively, mask fields if confidential tags presented.
func (sm *SM) __parse(field string, v reflect.Value, tag string) reflect.Value {
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
	// first iteration always get to Struct case, otherwise tag value always valid.
	switch orig.Kind() {
	case reflect.Struct:
		// copy type of original entity.
		t := orig.Type()
		// get new ptr of just created entity with original entity's type.
		cpVal = reflect.New(t)

		// iterate over struct fields
		for i := 0; i < t.NumField(); i++ {
			// get copied struct type's field i .
			f := t.Field(i)
			// get original struct type's field i.
			fVal := orig.Field(i)

			// if we won't be able to get Interface of field, skip it to prevent panic.
			if !fVal.CanInterface() {
				continue
			}

			// get confidential tag value
			tagVal := f.Tag.Get(sm.cfg.TagName)
			// if field type's kind is ptr or string, cast masked data to interface{} and then to string,
			// only after type casting set masked value to copied struct field.
			if fVal.Type().Kind() == reflect.Ptr && fVal.Elem().Kind() == reflect.String {
				s := sm.__parse(f.Name, fVal.Elem(), tagVal).Interface().(string)
				cpVal.Elem().Field(i).Set(reflect.ValueOf(&s))
			} else {
				cpVal.Elem().Field(i).Set(sm.__parse(f.Name, fVal, tagVal))
			}
		}
	case reflect.String:
		cpVal = reflect.New(orig.Type())
		cpVal.Elem().SetString(sm.__handle(v.String(), tag))
	case reflect.Slice, reflect.Array:
		cpVal = reflect.MakeSlice(orig.Type(), orig.Len(), orig.Cap())

		// slice/array values could not have tags.
		for i := 0; i < orig.Len(); i++ {
			cpVal.Index(i).Set(sm.__parse(field, orig.Index(i), ""))
		}
	case reflect.Map:
		cpVal = reflect.MakeMap(orig.Type())
		keys := orig.MapKeys()

		// map values could not have tags.
		for i := 0; i < orig.Len(); i++ {
			cpVal.SetMapIndex(keys[i], sm.__parse(keys[i].String(), orig.MapIndex(keys[i]), ""))
		}
	case reflect.Interface:
		cpVal = reflect.New(orig.Type())

		// try to mask data if interface{}'s underlying value is string.
		s, ok := v.Interface().(string)
		if ok {
			cpVal.Elem().Set(reflect.ValueOf(sm.__handle(s, tag)))
		} else {
			cpVal.Elem().Set(orig)
		}
	default:
		// just copy original values to copied struct type.
		cpVal = reflect.New(orig.Type())
		cpVal.Elem().Set(orig)
	}

	return sm.__adjust(v.Kind(), cpVal)
}

// __adjust returns correct reflect values according to its kind.
func (sm *SM) __adjust(k reflect.Kind, v reflect.Value) reflect.Value {
	switch k {
	case reflect.Map, reflect.Slice, reflect.Array, reflect.Ptr:
		return v
	default:
		return v.Elem()
	}
}

// __handle returns string parsed with Handler according to tag.
// If tag is empty, returns initial string.
func (sm *SM) __handle(input string, tag string) string {
	if handler := StructMasker.getHandler(tag); handler != nil {
		return handler(input)
	}

	return input
}
