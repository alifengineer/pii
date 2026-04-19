package mask

import "reflect"

func Mask(i any) any {
	rv := reflect.ValueOf(i)

	switch rv.Kind() {
	case reflect.Struct:
		return maskStruct(rv).Elem().Interface()
	case reflect.Pointer:
		if rv.Elem().Kind() == reflect.Struct {
			return maskStruct(rv.Elem()).Elem().Interface()
		}
	}
	return i
}

func maskStruct(rv reflect.Value) reflect.Value {
	rt := rv.Type()
	newRv := reflect.New(rt)

	for field := range rt.Fields() {
		fieldValue := rv.FieldByName(field.Name)
		var newFieldValue reflect.Value

		// if field has `pii` tag, masking field value
		if value, ok := field.Tag.Lookup(piiTag); ok {
			switch fieldValue.Kind() {
			case reflect.Pointer:
				if fieldValue.IsNil() {
					continue
				}
				newFieldValue = maskPtr(fieldValue, value)
			case reflect.Struct:
				newFieldValue = maskStruct(fieldValue).Elem()
			default:
				newFieldValue = maskLiterals(fieldValue, value)
			}
		} else {
			switch fieldValue.Kind() {
			case reflect.Struct:
				newFieldValue = maskStruct(fieldValue).Elem()
			default:
				newFieldValue = fieldValue
			}
		}

		newRv.Elem().FieldByName(field.Name).Set(newFieldValue)
	}

	return newRv
}

func maskPtr(fieldValue reflect.Value, tagVal string) reflect.Value {
	switch fieldValue.Kind() {
	case reflect.Struct:
		return maskStruct(fieldValue).Elem()
	case reflect.Pointer:
		if fieldValue.IsNil() {
			return fieldValue
		}

		newVal := maskPtr(fieldValue.Elem(), tagVal)
		ptr := reflect.New(newVal.Type())
		ptr.Elem().Set(newVal)

		return ptr
	}

	return maskLiterals(fieldValue, tagVal)
}

func maskLiterals(fv reflect.Value, tagVal string) reflect.Value {
	m, ok := maskers[TagValue(tagVal)]
	if !ok || fv.Kind() != reflect.String {
		return fv
	}

	masked := m(fv.String())
	result := reflect.New(fv.Type()).Elem()
	result.SetString(masked)
	return result
}
