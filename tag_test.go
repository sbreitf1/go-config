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
	NonZeroPrint     interface{} `config:"print:[nonzero]"`
	LenPrint         interface{} `config:"print:[len]"`
	MaskedPrint      interface{} `config:"print:[mask]"`
	HashedPrint      interface{} `config:"print:[sha256]"`
	NonZeroPrintName interface{} `config:"print:OtherName:[nonzero]"`
	LenPrintName     interface{} `config:"print:OtherName:[len]"`
	MaskedPrintName  interface{} `config:"print:OtherName:[mask]"`
	HashedPrintName  interface{} `config:"print:OtherName:[sha256]"`
	PrintName        interface{} `config:"print:VisibleName"`
	Default          interface{} `config:"default:some str"`
	DefaultWithColon interface{} `config:"default:some:nice:str"`
	FieldName        interface{} `config:"env:Bar,print:Bar,name:Foo"`
}

type tagTestCase struct {
	FieldName   string
	ExpectedTag tag
}

var tagTestCases = []tagTestCase{
	{"None", tag{"None", false, printModeDefault, "None", "None", "None", "", false}},
	{"Empty", tag{"Empty", false, printModeDefault, "Empty", "Empty", "Empty", "", false}},
	{"Flags", tag{"Flags", true, printModeDefault, "Flags", "Flags", "Flags", "", false}},
	{"EnvName", tag{"EnvName", false, printModeDefault, "EnvName", "SomeNewName", "EnvName", "", false}},
	{"NoPrint", tag{"NoPrint", false, printModeNone, "", "NoPrint", "NoPrint", "", false}},
	{"NonZeroPrint", tag{"NonZeroPrint", false, printModeNonZero, "NonZeroPrint", "NonZeroPrint", "NonZeroPrint", "", false}},
	{"LenPrint", tag{"LenPrint", false, printModeLen, "LenPrint", "LenPrint", "LenPrint", "", false}},
	{"MaskedPrint", tag{"MaskedPrint", false, printModeMasked, "MaskedPrint", "MaskedPrint", "MaskedPrint", "", false}},
	{"HashedPrint", tag{"HashedPrint", false, printModeSHA256, "HashedPrint", "HashedPrint", "HashedPrint", "", false}},
	{"NonZeroPrintName", tag{"NonZeroPrintName", false, printModeNonZero, "OtherName", "NonZeroPrintName", "NonZeroPrintName", "", false}},
	{"LenPrintName", tag{"LenPrintName", false, printModeLen, "OtherName", "LenPrintName", "LenPrintName", "", false}},
	{"MaskedPrintName", tag{"MaskedPrintName", false, printModeMasked, "OtherName", "MaskedPrintName", "MaskedPrintName", "", false}},
	{"HashedPrintName", tag{"HashedPrintName", false, printModeSHA256, "OtherName", "HashedPrintName", "HashedPrintName", "", false}},
	{"PrintName", tag{"PrintName", false, printModeDefault, "VisibleName", "PrintName", "PrintName", "", false}},
	{"Default", tag{"Default", false, printModeDefault, "Default", "Default", "Default", "some str", true}},
	{"DefaultWithColon", tag{"DefaultWithColon", false, printModeDefault, "DefaultWithColon", "DefaultWithColon", "DefaultWithColon", "some:nice:str", true}},
	{"FieldName", tag{"Foo", false, printModeDefault, "Bar", "Bar", "Foo", "", false}},
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
