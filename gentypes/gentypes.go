package gentypes

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/spf13/cobra"
)

type AdditionalField interface {
	GetName() string
	GetType() string
	IsReadOnly() bool
}

type Config struct {
	FilePath                   string
	CollectionAdditionalFields map[string][]AdditionalField
	PrintSelectOptions         bool
}

func Register(app *pocketbase.PocketBase, cfg Config) {
	app.RootCmd.AddCommand(&cobra.Command{
		Use: "gen-types",
		Run: func(cmd *cobra.Command, args []string) {
			err := cfg.generateTypes(app)
			if err != nil {
				fmt.Printf("error: %v\n", err)
			}
		},
	})

	app.OnCollectionAfterUpdateSuccess().BindFunc(func(e *core.CollectionEvent) error {
		err := cfg.generateTypes(app)
		if err != nil {
			return err
		}

		return e.Next()
	})
}

func (c *Config) generateTypes(app *pocketbase.PocketBase) error {
	collections, err := app.FindAllCollections()
	if err != nil {
		return err
	}

	root, err := projectRoot()
	if err != nil {
		panic(err)
	}

	constsPath := filepath.Join(strings.Trim(root, "\t\n\r "), c.FilePath, "pocketbase.ts")
	fConst, err := os.Create(constsPath)
	if err != nil {
		return err
	}
	defer fConst.Close()
	fmt.Fprintf(fConst, "/* This file was automatically generated, changes will be overwritten. */\n\n")
	fmt.Fprintf(fConst, "import PocketBase, { RecordService } from \"pocketbase\";\n\n")
	c.printCollectionConstants(fConst, collections)
	c.printCollectionRecordMap(fConst, collections)
	printTypedPocketBase(fConst)

	outPath := filepath.Join(strings.Trim(root, "\t\n\r "), c.FilePath, "base.d.ts")
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString("/* This file was automatically generated, changes will be overwritten. */\n\n")

	c.printBaseType(f)

	if c.PrintSelectOptions {
		optionsPath := filepath.Join(strings.Trim(root, "\t\n\r "), c.FilePath, "select-options.d.ts")

		fOptions, err := os.Create(optionsPath)
		if err != nil {
			return err
		}
		defer fOptions.Close()

		fOptions.WriteString("/* This file was automatically generated, changes will be overwritten. */\n\n")

		for _, collection := range collections {
			if !collection.System {
				if c.PrintSelectOptions {
					c.printCollectionSelectOptions(fOptions, collection)
				}
			}
		}

	}

	for _, collection := range collections {
		if !collection.System {
			c.printCollectionTypes(f, collection)
		}
	}

	return nil
}
