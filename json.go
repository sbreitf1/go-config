package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
)

// FromFile reads a JSON file and updates the given configuration.
//
// Respects the default json tag values.
func FromFile(path string, conf interface{}) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return FromJSON(data, conf)
}

// FromJSON parses JSON data and updates the given configuration.
//
// Respects the default json tag values.
func FromJSON(data []byte, conf interface{}) error {
	var obj interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	dst := newObject(conf)
	if !dst.IsAssignable() {
		return fmt.Errorf("conf must be an assignable value")
	}
	return fromJSON(obj, newPathPrefix(""), dst, nil)
}

func fromJSON(obj interface{}, prefix pathPrefix, dst *object, tag *tag) error {
	//TODO custom types with interfaces

	if dst.Is(typeDateTime) {
		return dateTimeFromJSON(obj, prefix, dst, tag)
	}
	if dst.Is(typeDuration) {
		return durationFromJSON(obj, prefix, dst, tag)
	}

	switch dst.Kind() {
	case reflect.Ptr:
		return fromJSON(obj, prefix, dst.Elem(), tag)

	case reflect.Struct:
		return structFromJSON(obj, prefix, dst)

	/*case reflect.Slice:
		return sliceFromEnvironment(prefix, dst)
	case reflect.Array:
		return arrayFromEnvironment(prefix, dst)*/

	case reflect.String:
		return stringFromJSON(obj, prefix, dst, tag)
	case reflect.Bool:
		return boolFromJSON(obj, prefix, dst, tag)
	case reflect.Int:
		return intFromJSON(obj, prefix, dst, tag)

	default:
		// just ignore unsupported types
		return nil
	}
}

func structFromJSON(obj interface{}, prefix pathPrefix, dst *object) error {
	t := reflect.TypeOf(obj)
	if t.Kind() != reflect.Map || t.Key().Kind() != reflect.String {
		return fmt.Errorf("%s: cannot parse struct from type %T", prefix.String(), obj)
	}

	// source obj must be a map
	src := make(map[string]interface{})
	v := reflect.ValueOf(obj)
	for _, key := range v.MapKeys() {
		src[key.Interface().(string)] = v.MapIndex(key).Interface()
	}

	//TODO use dst.IterateStruct
	fieldCount := dst.t.NumField()
	for i := 0; i < fieldCount; i++ {
		field := dst.t.Field(i)
		tag := getTag(field)
		val := dst.v.Field(i)

		fieldName := tag.FieldName
		// field name can be overwritten by json tag
		if jsonTag := field.Tag.Get("json"); len(jsonTag) > 0 {
			parts := strings.Split(jsonTag, ",")
			if len(parts) == 1 && parts[0] == "-" {
				// do not allow json input for this field
				continue
			}
			fieldName = parts[0]
		}

		if obj, ok := src[fieldName]; ok {
			if err := fromJSON(obj, prefix.Field(fieldName), &object{val.Type(), val}, &tag); err != nil {
				return err
			}
		}
	}
	return nil
}

/*func sliceFromEnvironment(prefix pathPrefix, dst *object) error {
	numStrVal, ok := lookupEnv(prefix.Field("Num").Env())
	if !ok || len(numStrVal) == 0 {
		return nil
	}

	num, err := strconv.Atoi(numStrVal)
	if err != nil {
		return fmt.Errorf("%s: failed to parse list length from %q", prefix.String(), numStrVal)
	}

	dst.InitSlice(num)
	return dst.IterateSlice(func(i int, dst *object) error {
		return fromEnvironment(prefix.Index(i), dst, nil)
	})
}

func arrayFromEnvironment(prefix pathPrefix, dst *object) error {
	return dst.IterateArray(func(i int, dst *object) error {
		return fromEnvironment(prefix.Index(i), dst, nil)
	})
}
*/

func stringFromJSON(obj interface{}, prefix pathPrefix, dst *object, tag *tag) error {
	return dst.SetString(obj.(string))
}

func boolFromJSON(obj interface{}, prefix pathPrefix, dst *object, tag *tag) error {
	return dst.SetBool(obj.(bool))
}

func intFromJSON(obj interface{}, prefix pathPrefix, dst *object, tag *tag) error {
	// json.Unmarshal treats all numbers as float64
	return dst.SetInt(int(obj.(float64)))
}

func dateTimeFromJSON(obj interface{}, prefix pathPrefix, dst *object, tag *tag) error {
	return dst.SetDateTimeFromString(obj.(string))
}

func durationFromJSON(obj interface{}, prefix pathPrefix, dst *object, tag *tag) error {
	return dst.SetDurationFromString(obj.(string))
}
