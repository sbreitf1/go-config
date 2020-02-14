package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const (
	envSeparator = "_"
)

var (
	lookupEnv = os.LookupEnv
	boolMap   = make(map[string]bool)
)

func init() {
	boolMap["true"] = true
	boolMap["yes"] = true
	boolMap["on"] = true
	boolMap["t"] = true
	boolMap["y"] = true
	boolMap["1"] = true
	boolMap["false"] = false
	boolMap["no"] = false
	boolMap["off"] = false
	boolMap["f"] = false
	boolMap["n"] = false
	boolMap["0"] = false
}

// FromEnvironment reads all values from environment variables.
func FromEnvironment(prefix string, conf interface{}) error {
	t := reflect.TypeOf(conf)
	v := reflect.ValueOf(conf)
	if !v.CanSet() && !(v.Kind() == reflect.Ptr && v.Elem().CanSet()) {
		return fmt.Errorf("conf must be an assignable value")
	}
	return fromEnvironment(strings.ToUpper(prefix), t, v, nil)
}

func fromEnvironment(prefix string, dstType reflect.Type, dstValue reflect.Value, tag *tag) error {
	//TODO advanced types like time.Time and time.Duration
	//TODO custom types with interfaces

	switch dstType.Kind() {
	case reflect.Ptr:
		return fromEnvironment(prefix, dstType.Elem(), dstValue.Elem(), tag)

	case reflect.Struct:
		return structFromEnvironment(prefix, dstType, dstValue)

	case reflect.Slice:
		return sliceFromEnvironment(prefix, dstType, dstValue)
	case reflect.Array:
		return arrayFromEnvironment(prefix, dstType, dstValue)

	case reflect.String:
		return stringFromEnvironment(prefix, dstType, dstValue, tag)
	case reflect.Bool:
		return boolFromEnvironment(prefix, dstType, dstValue, tag)
	case reflect.Int:
		return intFromEnvironment(prefix, dstType, dstValue, tag)

	default:
		// just ignore unsupported types
		return nil
	}
}

func structFromEnvironment(prefix string, dstType reflect.Type, dstValue reflect.Value) error {
	fieldCount := dstType.NumField()
	for i := 0; i < fieldCount; i++ {
		field := dstType.Field(i)
		tag := getTag(field)
		val := dstValue.Field(i)

		if val.CanSet() {
			if err := fromEnvironment(prefix+envSeparator+strings.ToUpper(tag.EnvName), field.Type, dstValue.Field(i), &tag); err != nil {
				return err
			}
		}
	}
	return nil
}

func sliceFromEnvironment(prefix string, dstType reflect.Type, dstValue reflect.Value) error {
	numStrVal, ok := lookupEnv(prefix + envSeparator + "NUM")
	if !ok || len(numStrVal) == 0 {
		return nil
	}

	num, err := strconv.Atoi(numStrVal)
	if err != nil {
		return fmt.Errorf("failed to parse list length for %s", prefix)
	}

	dstValue.Set(reflect.MakeSlice(dstType, num, num))
	for i := 0; i < num; i++ {
		fromEnvironment(prefix+envSeparator+strconv.Itoa(i), dstValue.Index(i).Type(), dstValue.Index(i), nil)
	}

	return nil
}

func arrayFromEnvironment(prefix string, dstType reflect.Type, dstValue reflect.Value) error {
	len := dstValue.Len()
	for i := 0; i < len; i++ {
		fromEnvironment(prefix+envSeparator+strconv.Itoa(i), dstValue.Index(i).Type(), dstValue.Index(i), nil)
	}

	return nil
}

func stringFromEnvironment(prefix string, dstType reflect.Type, dstValue reflect.Value, tag *tag) error {
	if strVal, ok := fromEnvOrDefault(prefix, tag); ok {
		dstValue.SetString(strVal)
	}
	return nil
}

func boolFromEnvironment(prefix string, dstType reflect.Type, dstValue reflect.Value, tag *tag) error {
	if strVal, ok := fromEnvOrDefault(prefix, tag); ok {
		val, ok := boolMap[strVal]
		if !ok {
			return fmt.Errorf("cannot parse bool from %q", strVal)
		}
		dstValue.SetBool(val)
	}
	return nil
}

func intFromEnvironment(prefix string, dstType reflect.Type, dstValue reflect.Value, tag *tag) error {
	if strVal, ok := fromEnvOrDefault(prefix, tag); ok {
		val, err := strconv.Atoi(strVal)
		if err != nil {
			return err
		}
		dstValue.SetInt(int64(val))
	}
	return nil
}

func fromEnvOrDefault(key string, tag *tag) (string, bool) {
	// explicit configuration from environment has highest priority
	if strVal, ok := lookupEnv(key); ok {
		return strVal, true
	}
	// no env available? try default value
	if tag.HasDefault {
		return tag.Default, true
	}
	// is not configured at all
	return "", false
}
