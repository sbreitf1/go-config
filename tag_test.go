package config

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TagTest struct {
	None             interface{}
	Empty            interface{} `config:""`
	Flags            interface{} `config:"required"`
	EnvName          interface{} `config:"env:SomeNewName"`
	NoPrint          interface{} `config:"print:-"`
	MaskedPrint      interface{} `config:"print:[mask]"`
	HashedPrint      interface{} `config:"print:[sha256]"`
	MaskedPrintName  interface{} `config:"print:OtherName:[mask]"`
	HashedPrintName  interface{} `config:"print:OtherName:[sha256]"`
	PrintName        interface{} `config:"print:VisibleName"`
	Default          interface{} `config:"default:some str"`
	DefaultWithColon interface{} `config:"default:some:nice:str"`
}

type tagTestCase struct {
	FieldName   string
	ExpectedTag tag
}

var tagTestCases = []tagTestCase{
	tagTestCase{"None", tag{"None", false, printModeDefault, "None", "None", "", false}},
	tagTestCase{"Empty", tag{"Empty", false, printModeDefault, "Empty", "Empty", "", false}},
	tagTestCase{"Flags", tag{"Flags", true, printModeDefault, "Flags", "Flags", "", false}},
	tagTestCase{"EnvName", tag{"EnvName", false, printModeDefault, "EnvName", "SomeNewName", "", false}},
	tagTestCase{"NoPrint", tag{"NoPrint", false, printModeNone, "", "NoPrint", "", false}},
	tagTestCase{"MaskedPrint", tag{"MaskedPrint", false, printModeMasked, "MaskedPrint", "MaskedPrint", "", false}},
	tagTestCase{"HashedPrint", tag{"HashedPrint", false, printModeSHA256, "HashedPrint", "HashedPrint", "", false}},
	tagTestCase{"MaskedPrintName", tag{"MaskedPrintName", false, printModeMasked, "OtherName", "MaskedPrintName", "", false}},
	tagTestCase{"HashedPrintName", tag{"HashedPrintName", false, printModeSHA256, "OtherName", "HashedPrintName", "", false}},
	tagTestCase{"PrintName", tag{"PrintName", false, printModeDefault, "VisibleName", "PrintName", "", false}},
	tagTestCase{"Default", tag{"Default", false, printModeDefault, "Default", "Default", "some str", true}},
	tagTestCase{"DefaultWithColon", tag{"DefaultWithColon", false, printModeDefault, "DefaultWithColon", "DefaultWithColon", "some:nice:str", true}},
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
