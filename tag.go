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
		options := strings.Split(tagStr, ",")
		for _, part := range options {
			switch part {
			case "required":
				tag.Required = true

			default:
				parts := strings.Split(part, ":")

				if len(parts) >= 2 {
					switch parts[0] {
					case "env":
						tag.EnvName = parts[1]
						continue

					case "print":
						if setPrintOptions(parts, &tag) {
							continue
						}

					case "default":
						tag.Default = strings.Join(parts[1:], ":")
						tag.HasDefault = true
						continue
					}
				}

				panic(fmt.Sprintf("cannot parse config tag from %q", tagStr))
			}
		}
	}

	return tag
}

func setPrintOptions(parts []string, tag *tag) bool {
	if len(parts) == 2 {
		switch parts[1] {
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
			tag.PrintName = parts[1]
		}
		return true
	}

	if len(parts) == 3 {
		tag.PrintName = parts[1]
		switch parts[2] {
		case "[len]":
			tag.PrintMode = printModeLen
			return true

		case "[mask]":
			tag.PrintMode = printModeMasked
			return true

		case "[sha256]":
			tag.PrintMode = printModeSHA256
			return true
		}
	}

	return false
}
