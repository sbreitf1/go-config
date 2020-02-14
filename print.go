package config

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	printModeDefault printMode = ""
	printModeNone    printMode = "none"
	printModeMasked  printMode = "masked"
	printModeSHA256  printMode = "sha256"
)

type printMode string

type printLine struct {
	Key   string
	Value interface{}
	Tag   *tag
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
	sprint(&lines, prefix, reflect.TypeOf(conf), reflect.ValueOf(conf), nil)

	maxKeyLen := 0
	for _, l := range lines {
		if len(l.Key) > maxKeyLen {
			maxKeyLen = len(l.Key)
		}
	}

	strLines := make([]string, len(lines))
	for i, l := range lines {
		//TODO respect print mode
		strLines[i] = fmt.Sprintf("%s:%s%v", l.Key, strings.Repeat(" ", maxKeyLen-len(l.Key)+1), l.Value)
	}
	return strLines
}

func sprint(lines *[]printLine, prefix string, dstType reflect.Type, dstValue reflect.Value, tag *tag) {
	//TODO correctly handle print mode inheritance from tag

	switch dstType.Kind() {
	case reflect.Ptr:
		sprint(lines, prefix, dstType.Elem(), dstValue.Elem(), tag)

	case reflect.Struct:
		sprintStruct(lines, prefix, dstType, dstValue, nil)

	default:
		*lines = append(*lines, printLine{prefix, dstValue.Interface(), tag})
	}
}

func sprintStruct(lines *[]printLine, prefix string, dstType reflect.Type, dstValue reflect.Value, tag *tag) {
	fieldCount := dstType.NumField()
	for i := 0; i < fieldCount; i++ {
		field := dstType.Field(i)
		tag := getTag(field)
		val := dstValue.Field(i)

		if val.CanInterface() {
			if tag.PrintMode != printModeNone {
				sprint(lines, prefix+"."+tag.PrintName, field.Type, val, &tag)
			}
		}
	}
}
