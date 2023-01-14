package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	validateCmd = &cobra.Command{
		Use:   "validate",
		Short: "Validates a Warden file to match the schema",
		RunE: func(cmd *cobra.Command, args []string) error {

			rawRule, err := loadWardenFile(wardenFileFl)
			if err != nil {
				log.Fatal(err)
			}

			schemaFile, err := os.Open("schema.json")
			if err != nil {
				log.Fatal(err)
			}

			var m interface{}

			err = yaml.Unmarshal(rawRule, &m)
			if err != nil {
				log.Fatal(err)
			}

			compiler := jsonschema.NewCompiler()
			if err := compiler.AddResource("schema.json", schemaFile); err != nil {
				log.Fatal(err)
			}

			schema, err := compiler.Compile("schema.json")
			if err != nil {
				log.Fatal(err)
			}

			if err := schema.Validate(m); err != nil {
				log.Fatal(err)
			}

			fmt.Println("Validation successful.")

			return nil
		},
	}
)

func init() {

	AddWardenFileFlag(validateCmd)

	rootCmd.AddCommand(validateCmd)
}
