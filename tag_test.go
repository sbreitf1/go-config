package config

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TagTest struct {
	None            interface{}
	Empty           interface{} `config:""`
	Flags           interface{} `config:"required"`
	EnvName         interface{} `config:"env:SomeNewName"`
	NoPrint         interface{} `config:"print:-"`
	MaskedPrint     interface{} `config:"print:[mask]"`
	HashedPrint     interface{} `config:"print:[sha256]"`
	MaskedPrintName interface{} `config:"print:OtherName:[mask]"`
	HashedPrintName interface{} `config:"print:OtherName:[sha256]"`
	PrintName       interface{} `config:"print:VisibleName"`
}

type tagTestCase struct {
	FieldName   string
	ExpectedTag tag
}

var tagTestCases = []tagTestCase{
	tagTestCase{"None", tag{"None", false, printModeDefault, "None", "None"}},
	tagTestCase{"Empty", tag{"Empty", false, printModeDefault, "Empty", "Empty"}},
	tagTestCase{"Flags", tag{"Flags", true, printModeDefault, "Flags", "Flags"}},
	tagTestCase{"EnvName", tag{"EnvName", false, printModeDefault, "EnvName", "SomeNewName"}},
	tagTestCase{"NoPrint", tag{"NoPrint", false, printModeNone, "", "NoPrint"}},
	tagTestCase{"MaskedPrint", tag{"MaskedPrint", false, printModeMasked, "MaskedPrint", "MaskedPrint"}},
	tagTestCase{"HashedPrint", tag{"HashedPrint", false, printModeSHA256, "HashedPrint", "HashedPrint"}},
	tagTestCase{"MaskedPrintName", tag{"MaskedPrintName", false, printModeMasked, "OtherName", "MaskedPrintName"}},
	tagTestCase{"HashedPrintName", tag{"HashedPrintName", false, printModeSHA256, "OtherName", "HashedPrintName"}},
	tagTestCase{"PrintName", tag{"PrintName", false, printModeDefault, "VisibleName", "PrintName"}},
}

func TestTags(t *testing.T) {
	for _, testCase := range tagTestCases {
		if !t.Run("TestTag"+testCase.FieldName, func(t *testing.T) {
			tag := getTagForField(testCase.FieldName)
			assert.Equal(t, testCase.ExpectedTag, tag)
		}) {
			break
		}
	}
}

func getTagForField(fieldName string) tag {
	var obj TagTest
	t := reflect.TypeOf(obj)
	field, ok := t.FieldByName(fieldName)
	if !ok {
		panic(fmt.Sprintf("field named %q not found in TagTest", fieldName))
	}
	return getTag(field)
}
