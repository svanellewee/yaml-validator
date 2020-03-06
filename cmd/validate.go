/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v2"
)

var YamlPath, SchemaPath string

func convert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = convert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = convert(v)
		}
	}
	return i
}

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		file, err := os.Open(YamlPath)
		if err != nil {
			log.Fatal(err)
		}
		decoder := yaml.NewDecoder(file)
		defer file.Close()
		var body interface{}
		if err := decoder.Decode(&body); err != nil {
			log.Fatal(err)
		}
		body = convert(body)

		jsonExample, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}
		releaseExample := gojsonschema.NewBytesLoader(jsonExample)

		schema := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s", SchemaPath))
		result, err := gojsonschema.Validate(schema, releaseExample)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("....", result.Errors(), result.Valid())
	},
}

func init() {
	validateCmd.Flags().StringVarP(&SchemaPath, "jsonschema", "s", "", "JSONSchema to validate against")
	validateCmd.Flags().StringVarP(&YamlPath, "input", "i", "", "input yaml to validate against")
	rootCmd.AddCommand(validateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// validateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// validateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
