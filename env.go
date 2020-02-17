package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"
)

var (
	lookupEnv = os.LookupEnv

	typeDateTime = reflect.TypeOf(time.Time{})
	typeDuration = reflect.TypeOf(time.Duration(0))
)

// FromEnvironment reads all values from environment variables.
func FromEnvironment(prefix string, conf interface{}) error {
	dst := newObject(conf)
	if !dst.IsAssignable() {
		return fmt.Errorf("conf must be an assignable value")
	}
	return fromEnvironment(newPathPrefix(prefix), dst, nil)
}

func fromEnvironment(prefix pathPrefix, dst *object, tag *tag) error {
	//TODO advanced types like time.Time and time.Duration
	//TODO custom types with interfaces

	if dst.Is(typeDateTime) {
		return dateTimeFromEnvironment(prefix, dst, tag)
	}
	if dst.Is(typeDuration) {
		return durationFromEnvironment(prefix, dst, tag)
	}

	switch dst.Kind() {
	case reflect.Ptr:
		return fromEnvironment(prefix, dst.Elem(), tag)

	case reflect.Struct:
		return structFromEnvironment(prefix, dst)

	case reflect.Slice:
		return sliceFromEnvironment(prefix, dst)
	case reflect.Array:
		return arrayFromEnvironment(prefix, dst)

	case reflect.String:
		return stringFromEnvironment(prefix, dst, tag)
	case reflect.Bool:
		return boolFromEnvironment(prefix, dst, tag)
	case reflect.Int:
		return intFromEnvironment(prefix, dst, tag)

	default:
		// just ignore unsupported types
		return nil
	}
}

func structFromEnvironment(prefix pathPrefix, dst *object) error {
	return dst.IterateStruct(func(dst *object, tag tag) error {
		if dst.IsAssignable() {
			return fromEnvironment(prefix.Field2(tag.FieldName, tag.EnvName), dst, &tag)
		}
		return nil
	})
}

func sliceFromEnvironment(prefix pathPrefix, dst *object) error {
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

func stringFromEnvironment(prefix pathPrefix, dst *object, tag *tag) error {
	return assignFromEnvOrDefault(prefix, dst.SetString, tag)
}

func boolFromEnvironment(prefix pathPrefix, dst *object, tag *tag) error {
	return assignFromEnvOrDefault(prefix, dst.SetBoolFromString, tag)
}

func intFromEnvironment(prefix pathPrefix, dst *object, tag *tag) error {
	return assignFromEnvOrDefault(prefix, dst.SetIntFromString, tag)
}

func dateTimeFromEnvironment(prefix pathPrefix, dst *object, tag *tag) error {
	return assignFromEnvOrDefault(prefix, dst.SetDateTimeFromString, tag)
}

func durationFromEnvironment(prefix pathPrefix, dst *object, tag *tag) error {
	return assignFromEnvOrDefault(prefix, dst.SetDurationFromString, tag)
}

func assignFromEnvOrDefault(prefix pathPrefix, assignHandler func(string) error, tag *tag) error {
	if strVal, ok := fromEnvOrDefault(prefix.Env(), tag); ok {
		if err := assignHandler(strVal); err != nil {
			return fmt.Errorf("%s: %s", prefix.String(), err.Error())
		}
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
