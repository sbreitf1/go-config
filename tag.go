package config

import (
	"fmt"
	"reflect"
	"strings"
)

type tag struct {
	FieldName  string
	Required   bool
	PrintMode  printMode
	PrintName  string
	EnvName    string
	Default    string
	HasDefault bool
}

func getTag(field reflect.StructField) tag {
	tag := tag{
		FieldName: field.Name,
		Required:  false,
		PrintMode: printModeDefault,
		PrintName: field.Name,
		EnvName:   field.Name,
	}

	tagStr := field.Tag.Get("config")
	if len(tagStr) > 0 {
		options := make(map[string][]string)
		for _, val := range strings.Split(tagStr, ",") {
			parts := strings.Split(val, ":")
			opt := parts[0]
			//TODO deny unknown config options
			if _, ok := options[opt]; ok {
				panic(fmt.Sprintf("config option %q specified multiple times", opt))
			}
			options[opt] = parts[1:]
		}

		if args, ok := options["required"]; ok {
			if len(args) != 0 {
				panic("config option \"required\" does not support any arguments")
			}
			tag.Required = true
		}

		if args, ok := options["name"]; ok {
			if len(args) != 1 {
				panic("config option \"name\" requires exactly one argument")
			}
			// needs to be evaluated before all other names to prevent overrides
			tag.FieldName = args[0]
			tag.EnvName = args[0]
			tag.PrintName = args[0]
		}

		if args, ok := options["env"]; ok {
			if len(args) != 1 {
				panic("config option \"env\" requires exactly one argument")
			}
			tag.EnvName = args[0]
		}

		if args, ok := options["print"]; ok {
			setPrintOptions(args, &tag)
		}

		if args, ok := options["default"]; ok {
			tag.Default = strings.Join(args, ":")
			tag.HasDefault = true
		}
	}

	return tag
}

func setPrintOptions(args []string, tag *tag) {
	if len(args) == 0 {
		panic("config option \"print\" requires at least one argument")
	}

	if len(args) == 1 {
		switch args[0] {
		case "-":
			tag.PrintMode = printModeNone
			tag.PrintName = ""

		case "[len]":
			tag.PrintMode = printModeLen

		case "[mask]":
			tag.PrintMode = printModeMasked

		case "[sha256]":
			tag.PrintMode = printModeSHA256

		default:
			tag.PrintName = args[0]
		}
		return
	}

	if len(args) == 2 {
		tag.PrintName = args[0]
		switch args[1] {
		case "[len]":
			tag.PrintMode = printModeLen
			return

		case "[mask]":
			tag.PrintMode = printModeMasked
			return

		case "[sha256]":
			tag.PrintMode = printModeSHA256
			return
		}
		panic("")
	}

	panic("too many arguments for config option \"print\"")
}
