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
	fmt.Fprint(f, "  expand?: { [key: string]: any };")
	fmt.Fprint(f, "}\n\n")
}

func (c *Config) printCollectionSchema(f *os.File, collections []*core.Collection) {
	fmt.Fprintf(f, `import { CollectionRecords } from "./pocketbase-types";

type FieldType =
  | "text"
  | "number"
  | "bool"
  | "relation"
  | "select"
  | "json"
  | "file"
  | "date"
  | "autodate"
  | "email";

type FieldDefinition = {
  type: FieldType;
  values?: string[];
};

export type FieldSchema = {
  [C in keyof CollectionRecords]: {
    [F in keyof CollectionRecords[C]]?: FieldDefinition;
  };
};

export const fieldSchema: FieldSchema = {
`)

	for _, collection := range collections {
		if collection.System {
			continue
		}

		fmt.Fprintf(f, "  %s: {\n", collection.Name)

		for _, field := range collection.Fields {
			if field.GetHidden() {
				continue
			}
			fmt.Fprintf(f, "    %s: { type: \"%s\"", field.GetName(), field.Type())

			if sf, ok := field.(*core.SelectField); ok {
				fmt.Fprint(f, ", values: [")
				for i, v := range sf.Values {
					fmt.Fprintf(f, "%q", v)
					if i < len(sf.Values)-1 {
						fmt.Fprint(f, ", ")
					}
				}

				fmt.Fprint(f, "]")
			}
			fmt.Fprint(f, " },\n")
		}
		fmt.Fprint(f, "  },\n")
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
	fmt.Fprintln(f, "  collection<K extends keyof CollectionRecords>(name: K): RecordService<CollectionRecords[K]>;")
	fmt.Fprintln(f, "  // fallback for dynamic types")
	fmt.Fprintln(f, "  collection<T>(name: string): RecordService<T>;")
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
