package gentypes

import (
	"fmt"
	"strings"

	"github.com/pocketbase/pocketbase/core"
)

func toTypeScriptType(f core.Field) string {
	switch f.Type() {
	case "password":
		return ": string"
	case "text":
		return ": string"
	case "email":
		return ": string"
	case "relation":
		if sf, ok := f.(*core.RelationField); ok {
			res := ""
			res += ": string"
			if sf.MaxSelect > 1 {
				res += "[]"
			}
			return res
		}
		return ": string"
	case "autodate":
		return ": string"
	case "date":
		return ": string"
	case "url":
		return ": string"
	case "file":
		return ": string"
	case "select":
		if sf, ok := f.(*core.SelectField); ok {
			res := ""
			values := sf.Values

			if !sf.Required {
				values = append(values, "")
			}

			if len(values) > 0 {
				var quoted []string
				for _, v := range values {
					quoted = append(quoted, fmt.Sprintf("\"%s\"", v))
				}
				res = strings.Join(quoted, " | ")
			}

			if sf.MaxSelect > 1 {
				res = fmt.Sprintf("(%s)[]", res)
			}
			return fmt.Sprintf(": %s", res)
		}
		return "?: string"
	case "number":
		return ": number"
	case "bool":
		return ": boolean"
	case "json":
		return ": any"
	default:
		return ": unknown"
	}
}

func additionalFieldToTypeScriptType(fType string) string {
	res := ""

	switch fType {
	case "text":
		res = ": string"
	case "number":
		res = ": number"
	case "bool":
		res = ": boolean"
	case "json":
		res = ": any"
	default:
		res = ": unknown"
	}

	return fmt.Sprintf("%s; // %s\n", res, fType)
}
