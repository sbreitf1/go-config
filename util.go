package config

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type pathPrefix []interface{}

type fieldName struct {
	RealName    string
	VisibleName string
}

func (p pathPrefix) String() string {
	var sb strings.Builder
	for i, pathPart := range p {
		switch p := pathPart.(type) {
		case fieldName:
			if i > 0 {
				sb.WriteString(".")
			}
			sb.WriteString(p.RealName)
		case int:
			sb.WriteString("[")
			sb.WriteString(strconv.Itoa(p))
			sb.WriteString("]")
		default:
			panic(fmt.Sprintf("invalid path part of type %T", pathPart))
		}
	}
	return sb.String()
}

func (p pathPrefix) Env() string {
	var sb strings.Builder
	for i, pathPart := range p {
		if i > 0 && sb.Len() > 0 {
			sb.WriteString("_")
		}
		switch p := pathPart.(type) {
		case fieldName:
			sb.WriteString(strings.ToUpper(p.VisibleName))
		case int:
			sb.WriteString(strconv.Itoa(p))
		default:
			panic(fmt.Sprintf("invalid path part of type %T", pathPart))
		}
	}
	return sb.String()
}

func (p pathPrefix) Field(name string) pathPrefix {
	return append(p, fieldName{name, name})
}

func (p pathPrefix) Field2(realName, visibleName string) pathPrefix {
	return append(p, fieldName{realName, visibleName})
}

func (p pathPrefix) Index(index int) pathPrefix {
	return append(p, index)
}

func newPathPrefix(firstField string) pathPrefix {
	return []interface{}{fieldName{firstField, firstField}}
}

var (
	boolMap = make(map[string]bool)
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

type object struct {
	t reflect.Type
	v reflect.Value
}

func (obj *object) IsAssignable() bool {
	return obj.v.CanSet() || (obj.v.Kind() == reflect.Ptr && obj.v.Elem().CanSet())
}

func (obj *object) IsReadable() bool {
	return obj.v.CanInterface()
}

func (obj *object) Kind() reflect.Kind {
	return obj.t.Kind()
}

func (obj *object) NumField() int {
	return obj.t.NumField()
}

func (obj *object) Len() int {
	return obj.v.Len()
}

func (obj *object) Elem() *object {
	return &object{obj.t.Elem(), obj.v.Elem()}
}

func (obj *object) Index(i int) *object {
	return &object{obj.v.Index(i).Type(), obj.v.Index(i)}
}

func (obj *object) Interface() interface{} {
	return obj.v.Interface()
}

func (obj *object) IterateArray(f func(i int, obj *object) error) error {
	len := obj.v.Len()
	for i := 0; i < len; i++ {
		if err := f(i, obj.Index(i)); err != nil {
			return err
		}
	}
	return nil
}

func (obj *object) IterateSlice(f func(i int, obj *object) error) error {
	return obj.IterateArray(f)
}

func (obj *object) IterateStruct(f func(obj *object, tag tag) error) error {
	fieldCount := obj.t.NumField()
	for i := 0; i < fieldCount; i++ {
		field := obj.t.Field(i)
		tag := getTag(field)
		val := obj.v.Field(i)

		if err := f(&object{val.Type(), val}, tag); err != nil {
			return err
		}
	}
	return nil
}

func (obj *object) InitSlice(len int) {
	obj.v.Set(reflect.MakeSlice(obj.t, len, len))
}

func (obj *object) SetString(val string) error {
	obj.v.SetString(val)
	return nil
}

func (obj *object) SetBoolFromString(strVal string) error {
	val, ok := boolMap[strVal]
	if !ok {
		return fmt.Errorf("cannot parse bool from %q", strVal)
	}
	obj.v.SetBool(val)
	return nil
}

func (obj *object) SetIntFromString(strVal string) error {
	val, err := strconv.Atoi(strVal)
	if err != nil {
		return fmt.Errorf("cannot parse int from %q", strVal)
	}
	obj.v.SetInt(int64(val))
	return nil
}

func (obj *object) SetDateTimeFromString(strVal string) error {
	dt, err := func() (time.Time, error) {
		strVal = strings.Replace(strVal, " ", "T", -1)
		switch len(strVal) {
		case 10:
			return time.ParseInLocation("2006-01-02", strVal, time.Local)
		case 19:
			return time.ParseInLocation("2006-01-02T15:04:05", strVal, time.Local)
		case 20:
			return time.ParseInLocation("2006-01-02T15:04:05Z", strVal, time.UTC)
		case 24:
			return time.ParseInLocation("2006-01-02T15:04:05-0700", strVal, time.UTC)
		case 25:
			return time.ParseInLocation("2006-01-02T15:04:05-0700", strVal[:22]+strVal[23:], time.UTC)
		}
		return time.Time{}, fmt.Errorf("invalid format")
	}()
	if err != nil {
		return fmt.Errorf("cannot parse datetime from %q", strVal)
	}

	obj.v.Set(reflect.ValueOf(dt))
	return nil
}

func newObject(obj interface{}) *object {
	return &object{reflect.TypeOf(obj), reflect.ValueOf(obj)}
}
