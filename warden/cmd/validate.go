package cmd

import (
	"bytes"
	_ "embed"
	"fmt"
	"log"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

//go:embed policy.schema.json
var schemaFile []byte

var (
	validateCmd = &cobra.Command{
		Use:   "validate",
		Short: "Validates a Warden file to match the schema",
		RunE: func(cmd *cobra.Command, args []string) error {

			_, policyFile, err := loadPolicyFile(policyFileFl)
			if err != nil {
				log.Fatal(err)
			}

			schemaReader := bytes.NewReader(schemaFile)
			if err != nil {
				log.Fatal(err)
			}

			var m interface{}

			err = yaml.Unmarshal(policyFile, &m)
			if err != nil {
				log.Fatal(err)
			}

			compiler := jsonschema.NewCompiler()
			if err := compiler.AddResource("schema.json", schemaReader); err != nil {
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

	AddPolicyFileFlag(validateCmd)

	rootCmd.AddCommand(validateCmd)
}
