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

func Print(prefix string, conf interface{}) error {
	str, err := ToString(prefix, conf)
	if err != nil {
		return err
	}
	_, err = fmt.Println(str)
	return err
}

func ToString(prefix string, conf interface{}) (string, error) {
	lines, err := ToLines(prefix, conf)
	if err != nil {
		return "", err
	}
	return strings.Join(lines, "\n"), nil
}

func ToLines(prefix string, conf interface{}) ([]string, error) {
	lines := make([]printLine, 0)
	if err := sprint(&lines, prefix, reflect.TypeOf(conf), reflect.ValueOf(conf), nil); err != nil {
		return nil, err
	}

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
	return strLines, nil
}

func sprint(lines *[]printLine, prefix string, dstType reflect.Type, dstValue reflect.Value, tag *tag) error {
	//TODO correctly handle print mode inheritance from tag

	switch dstType.Kind() {
	case reflect.Ptr:
		return sprint(lines, prefix, dstType.Elem(), dstValue.Elem(), tag)

	case reflect.Struct:
		return sprintStruct(lines, prefix, dstType, dstValue, nil)

	default:
		*lines = append(*lines, printLine{prefix, dstValue.Interface(), tag})
		return nil
		//return fmt.Errorf("type %s not supported", dstType.Kind())
	}
}

func sprintStruct(lines *[]printLine, prefix string, dstType reflect.Type, dstValue reflect.Value, tag *tag) error {
	fieldCount := dstType.NumField()
	for i := 0; i < fieldCount; i++ {
		field := dstType.Field(i)
		tag := getTag(field)

		if tag.PrintMode != printModeNone {
			if err := sprint(lines, prefix+"."+tag.PrintName, field.Type, dstValue.Field(i), &tag); err != nil {
				return err
			}
		}
	}
	return nil
}
