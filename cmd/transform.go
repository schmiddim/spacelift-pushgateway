package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"spacelift-pushgateway/api"
)

var filename string

var transformCmd = &cobra.Command{
	Use:   "transform --file=filename",
	Short: "Transform JSON values based on specified paths and separators",
	Long: `Transform command takes a JSON file and applies transformations to specific paths within the JSON.
It reads the JSON file, applies transformations based on the given paths and separators, and then outputs the transformed JSON.`,
	Run: func(cmd *cobra.Command, args []string) {
		jsonData := readJsonFile(filename)
		var (
			transformedJSON []byte
			err             error
		)

		for _, splits := range config.Json.ValueSplits {
			transformedJSON, err = api.TransformJsonValues(jsonData, splits.Path, splits.Separator)
			if err != nil {
				log.Fatalf("Error transforming JSON: %v", err)
			}
		}
		fmt.Println(string(transformedJSON))
	},
}

func init() {
	transformCmd.Flags().StringVar(&filename, "file", "", "Path to the JSON file that will be transformed.")
	err := transformCmd.MarkFlagRequired("file")
	if err != nil {
		log.Fatal(err)
	}
	rootCmd.AddCommand(transformCmd)
}
