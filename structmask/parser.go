package structmask

import (
	"reflect"
)

// __normalize returns correct reflect values according to its kind.
func (sm *SM) __normalize(k reflect.Kind, v reflect.Value) reflect.Value {
	switch k {
	case reflect.Map, reflect.Slice, reflect.Array, reflect.Ptr:
		return v
	default:
		return v.Elem()
	}
}

// __parse parse data recursively, modify fields data if custom tag labels presented.
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
			label := f.Tag.Get(sm.tag)
			// if field type's kind is ptr or string, cast modified data to interface{} and then to initial type,
			// only after type casting set modified value to copied struct field.
			if fVal.Type().Kind() == reflect.Ptr && fVal.Elem().Kind() == reflect.String {
				s := sm.__parse(f.Name, fVal.Elem(), label).Interface().(string)
				cpVal.Elem().Field(i).Set(reflect.ValueOf(&s))
			} else {
				cpVal.Elem().Field(i).Set(sm.__parse(f.Name, fVal, label))
			}
		}
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
	default:
		cpVal = reflect.New(orig.Type())

		if modVal, ok := sm.__handle(v.Interface(), tag); ok {
			cpVal.Elem().Set(reflect.ValueOf(modVal))
		} else {
			cpVal.Elem().Set(orig)
		}
	}

	return sm.__normalize(v.Kind(), cpVal)
}

// __handle returns field's value parsed with Handler according to tag label.
// If it is empty, returns initial value.
func (sm *SM) __handle(input any, tag string) (any, bool) {
	inputKind := reflect.TypeOf(input).Kind()

	if handler, ok := StructMasker.getHandler(tag); ok {
		result := handler(input)
		resultKind := reflect.TypeOf(result).Kind()

		if resultKind == inputKind {
			return result, true
		}
	}

	return input, false
}
