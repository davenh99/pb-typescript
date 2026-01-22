package gentypes

import (
	"fmt"
	"os"

	"github.com/pocketbase/pocketbase/core"
)

func (c *Config) printBaseType(f *os.File) {
	fmt.Fprint(f, "interface BaseRecord {\n")

	baseFields := []string{"id", "collectionName", "collectionId", "created", "updated"}

	for _, field := range baseFields {
		fmt.Fprintf(f, "  readonly %s: string;\n", field)
	}

	fmt.Fprint(f, "}\n\n")
}

func (c *Config) printCollectionSelectOptions(f *os.File, collection *core.Collection) {
	collectionName := capitalise(collection.Name)

	selectFields := make([]*core.SelectField, 0)

	for _, field := range collection.Fields {
		if field.Type() == "select" && !field.GetHidden() {
			if sf, ok := field.(*core.SelectField); ok {
				selectFields = append(selectFields, sf)
			}
		}
	}

	if len(selectFields) == 0 {
		return
	}

	fmt.Fprintf(f, "export const %sSelectOptions = {\n", collectionName)

	for _, sf := range selectFields {
		fieldName := sf.GetName()

		values := append([]string{}, sf.Values...)

		fmt.Fprintf(f, "  %s: [", fieldName)
		for i, v := range values {
			fmt.Fprintf(f, "%q", v)
			if i < len(values)-1 {
				fmt.Fprint(f, ", ")
			}
		}
		fmt.Fprintln(f, "],")
	}

	fmt.Fprint(f, "};\n\n")
}

func (c *Config) printCollectionTypes(f *os.File, collection *core.Collection) {
	collectionName := capitalise(collection.Name)

	fmt.Fprintf(f, "/* Collection type: %s */\n", collection.Type)
	fmt.Fprintf(f, "interface %s {\n", collectionName)

	for _, field := range collection.Fields {
		if field.Type() == "autodate" || field.GetName() == "id" || field.GetHidden() {
			continue
		}
		fmt.Fprintf(f, "  %s%s; // %s\n", field.GetName(), toTypeScriptType(field), field.Type())
	}

	for _, additionalField := range c.CollectionAdditionalFields[collection.Name] {
		readonly := ""
		if additionalField.IsReadOnly() {
			readonly = "readonly "
		}
		fmt.Fprintf(
			f,
			"  %s%s%s",
			readonly,
			additionalField.GetName(),
			additionalFieldToTypeScriptType(additionalField.GetType()),
		)
	}

	fmt.Fprintln(f, "}")

	recordName := collectionName + "Record"
	fmt.Fprintf(f, "type %s = %s & BaseRecord;\n", recordName, collectionName)
	fmt.Fprintf(f, "type %sUpdatePayload = Partial<%s>;\n\n", collectionName, recordName)
}

func printTypedPocketBase(f *os.File) {
	fmt.Fprintln(f, "export interface TypedPocketBase extends PocketBase {")
	fmt.Fprintln(f, "  collection<K extends keyof CollectionRecords>(")
	fmt.Fprintln(f, "    name: K")
	fmt.Fprintln(f, "  ): RecordService<CollectionRecords[K]>;")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "  // fallback for dynamic strings")
	fmt.Fprintln(f, "  collection(name: string): RecordService<any>;")
	fmt.Fprintln(f, "}")
	fmt.Fprintln(f, "")
}

func (c *Config) printCollectionConstants(f *os.File, collections []*core.Collection) {
	fmt.Fprintln(f, "export const Collections = {")
	for _, col := range collections {
		if col.System {
			continue
		}
		fmt.Fprintf(
			f,
			"  %s: %q,\n",
			capitalise(col.Name),
			col.Name,
		)
	}
	fmt.Fprintf(f, "} as const;\n\n")
}

func (c *Config) printCollectionRecordMap(f *os.File, collections []*core.Collection) {
	fmt.Fprintln(f, "export interface CollectionRecords {")
	for _, col := range collections {
		if col.System {
			continue
		}
		fmt.Fprintf(
			f,
			"  %s: %sRecord;\n",
			col.Name,
			capitalise(col.Name),
		)
	}
	fmt.Fprintf(f, "}\n\n")
}
