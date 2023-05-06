package confidential

import (
	"reflect"
	"sexreflection/pkg/mask"
	"sexreflection/pkg/utils"
)

type SM struct{}

func newSM() *SM {
	return &SM{}
}

// Proceed replace tagged fields of struct with correct masks
//
// Tag example: `confidential:"cvv"`.
// All currently supported confidential tags: cvv, cardnumber, cardholder.
//
// See confidential/const.go
func (sm *SM) Proceed(input any) any {
	if input == nil ||
		(reflect.ValueOf(input).Kind() != reflect.Struct &&
			reflect.ValueOf(input).Kind() != reflect.Ptr) {
		return nil
	}

	return __parse("", reflect.ValueOf(input), "").Interface()
}

// You don't need to deal with it, really
// __parse parse data recursively, mask fields if confidential tags presented.
func __parse(field string, v reflect.Value, tag string) reflect.Value {
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
			tagVal := f.Tag.Get(confidentialTagKey)
			// if field type's kind is ptr or string, cast masked data to interface{} and then to string,
			// only after type casting set masked value to copied struct field.
			if fVal.Type().Kind() == reflect.Ptr && fVal.Elem().Kind() == reflect.String {
				s := __parse(f.Name, fVal.Elem(), tagVal).Interface().(string)
				cpVal.Elem().Field(i).Set(reflect.ValueOf(&s))
			} else {
				cpVal.Elem().Field(i).Set(__parse(f.Name, fVal, tagVal))
			}
		}
	case reflect.String:
		cpVal = reflect.New(orig.Type())
		cpVal.Elem().SetString(__mask(v.String(), tag))
	case reflect.Slice, reflect.Array:
		cpVal = reflect.MakeSlice(orig.Type(), orig.Len(), orig.Cap())

		// slice/array values could not have tags.
		for i := 0; i < orig.Len(); i++ {
			cpVal.Index(i).Set(__parse(field, orig.Index(i), ""))
		}
	case reflect.Map:
		cpVal = reflect.MakeMap(orig.Type())
		keys := orig.MapKeys()

		// map values could not have tags.
		for i := 0; i < orig.Len(); i++ {
			cpVal.SetMapIndex(keys[i], __parse(keys[i].String(), orig.MapIndex(keys[i]), ""))
		}
	case reflect.Interface:
		cpVal = reflect.New(orig.Type())

		// try to mask data if interface{}'s underlying value is string.
		s, ok := v.Interface().(string)
		if ok {
			cpVal.Elem().Set(reflect.ValueOf(__mask(s, tag)))
		} else {
			cpVal.Elem().Set(orig)
		}
	default:
		// just copy original values to copied struct type.
		cpVal = reflect.New(orig.Type())
		cpVal.Elem().Set(orig)
	}

	return __adjust(v.Kind(), cpVal)
}

// __adjust returns correct reflect values according to its kind.
func __adjust(k reflect.Kind, v reflect.Value) reflect.Value {
	switch k {
	case reflect.Map, reflect.Slice, reflect.Array, reflect.Ptr:
		return v
	default:
		return v.Elem()
	}
}

// __mask returns correct mask according to tag.
// If tag is empty, just return initial string.
func __mask(s string, tag string) string {
	if utils.AnyOf(tag, tags...) {
		switch tag {
		case cvv:
			return mask.CVV()
		case cardNumber:
			return mask.CardNumber(s)
		case cardHolder:
			return mask.CardHolder()
		case def:
			return mask.Default()
		}
	}

	return s
}
