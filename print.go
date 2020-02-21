package config

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	printModeDefault printMode = ""
	printModeNone    printMode = "-"
	printModeNonZero printMode = "nonzero"
	printModeLen     printMode = "len"
	printModeMasked  printMode = "masked"
	printModeSHA256  printMode = "sha256"
	//TODO printModeKey to print out map keys as if they are actual struct fields
	//TODO printModeEscape to print out in "" with escape sequences
)

type printMode string

type printLine struct {
	Key   string
	Value interface{}
	Mode  printMode
	Tag   *tag
}

func (l printLine) PrintVisible() (string, bool) {
	switch l.Mode {
	case printModeDefault:
		return fmt.Sprintf("%v", l.Value), true

	case printModeNonZero:
		if reflect.ValueOf(l.Value).IsZero() {
			return "", false
		}
		return fmt.Sprintf("%v", l.Value), true

	case printModeLen:
		return strconv.Itoa(reflect.ValueOf(l.Value).Len()), true

	case printModeMasked:
		if reflect.ValueOf(l.Value).IsZero() {
			return "", false
		}
		return "******", true

	case printModeSHA256:
		panic("print mode sha256 not yet supported")

	default:
		panic(fmt.Sprintf("unsupported print mode %q", l.Mode))
	}
}

// Print prints the output of ToString with fmt.Println.
func Print(prefix string, conf interface{}) error {
	_, err := fmt.Println(ToString(prefix, conf))
	return err
}

// ToString returns the output of ToLines concatenated with "\n".
func ToString(prefix string, conf interface{}) string {
	return strings.Join(ToLines(prefix, conf), "\n")
}

// ToLines returns a line for every configuration value with equal indentation of all values.
func ToLines(prefix string, conf interface{}) []string {
	lines := make([]printLine, 0)
	sprint(&lines, newPathPrefix(prefix), newObject(conf), printModeDefault, nil)

	type entry struct {
		Key   string
		Value string
	}

	entries := make([]entry, 0)

	maxKeyLen := 0
	for _, l := range lines {
		if visibleString, ok := l.PrintVisible(); ok {
			if len(l.Key) > maxKeyLen {
				maxKeyLen = len(l.Key)
			}
			entries = append(entries, entry{l.Key, visibleString})
		}
	}

	strLines := make([]string, len(entries))
	for i, e := range entries {
		strLines[i] = fmt.Sprintf("%s:%s%s", e.Key, strings.Repeat(" ", maxKeyLen-len(e.Key)+1), e.Value)
	}
	return strLines
}

func sprint(lines *[]printLine, prefix pathPrefix, obj *object, mode printMode, tag *tag) {
	if mode == printModeLen {
		// do not print full hierarchy, only number of elements:
		*lines = append(*lines, printLine{prefix.String(), obj.Interface(), mode, tag})
		return
	}

	switch obj.Kind() {
	case reflect.Ptr:
		if !obj.IsNil() {
			sprint(lines, prefix, obj.Elem(), mode, tag)
		}

	case reflect.Struct:
		sprintStruct(lines, prefix, obj, mode)

	case reflect.Slice:
		sprintSlice(lines, prefix, obj, mode)
	case reflect.Array:
		sprintArray(lines, prefix, obj, mode)

	default:
		*lines = append(*lines, printLine{prefix.String(), obj.Interface(), mode, tag})
	}
}

func sprintStruct(lines *[]printLine, prefix pathPrefix, obj *object, mode printMode) {
	obj.IterateStruct(func(obj *object, tag tag) error {
		if obj.IsReadable() {
			if tag.PrintMode != printModeNone {
				newMode := mode
				if tag.PrintMode != printModeDefault {
					newMode = tag.PrintMode
				}
				sprint(lines, prefix.Field(tag.PrintName), obj, newMode, &tag)
			}
		}
		return nil
	})
}

func sprintSlice(lines *[]printLine, prefix pathPrefix, obj *object, mode printMode) {
	obj.IterateArray(func(i int, obj *object) error {
		sprint(lines, prefix.Index(i), obj, mode, nil)
		return nil
	})
}

func sprintArray(lines *[]printLine, prefix pathPrefix, obj *object, mode printMode) {
	obj.IterateArray(func(i int, obj *object) error {
		sprint(lines, prefix.Index(i), obj, mode, nil)
		return nil
	})
}
