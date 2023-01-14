package cmd

import (
	"bytes"
	_ "embed"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

//go:embed schema.json
var schemaFile []byte

var (
	validateCmd = &cobra.Command{
		Use:   "validate",
		Short: "Validates a Warden file to match the schema",
		RunE: func(cmd *cobra.Command, args []string) error {

			wardenFilePath := "warden."

			_, err := os.Stat("warden.yml")
			if os.IsNotExist(err) {

				_, err := os.Stat("warden.yaml")
				if os.IsNotExist(err) {
					return fmt.Errorf("A Warden file was not located at ./warden.yml")
				} else if err != nil {
					return fmt.Errorf("There was a problem reading the Warden file.")
				} else {
					wardenFilePath = wardenFilePath + "yaml"
				}
			} else if err != nil {
				return fmt.Errorf("There was a problem reading the Warden file.")
			} else {
				wardenFilePath = wardenFilePath + "yml"
			}

			rawRule, err := ioutil.ReadFile(wardenFilePath)
			if err != nil {
				log.Fatal(err)
			}

			schemaReader := bytes.NewReader(schemaFile)
			if err != nil {
				log.Fatal(err)
			}

			var m interface{}

			err = yaml.Unmarshal(rawRule, &m)
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
	rootCmd.AddCommand(validateCmd)
}
