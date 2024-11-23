package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"spacelift-pushgateway/api"
)

var transformBeforeExtract = true
var rename = true

var extractCmd = &cobra.Command{
	Use:   "extract --file=filename",
	Short: "Extracts specific fields from a JSON file",
	Long:  `The extract command reads a JSON file and extracts specified fields.`,
	Run: func(cmd *cobra.Command, args []string) {

		jsonData := readJsonFile(filename)
		var (
			err error
		)

		if transformBeforeExtract {
			for _, splits := range config.Json.ValueSplits {
				jsonData, err = api.TransformJsonValues(jsonData, splits.Path, splits.Separator)
				if err != nil {
					log.Fatalf("Error transforming JSON: %v", err)
				}
			}
		}
		results, err := api.ExtractMultipleJSONPaths(jsonData, config.Json.FieldsToExtract)
		if err != nil {
			log.Fatalf("Error extracting data: %v", err)
		}
		if rename {
			results, err = api.RenameKeys(results, config.Json.Rename)
		}

		gw := api.NewPushGateway(config.Prometheus.PushGatewayUrl, config.Prometheus.TargetMetric, config.Prometheus.TargetMetricHelp, config.Prometheus.JobName)
		hasErrors, errors := gw.ValidateLabels(results)
		if hasErrors == true {
			for _, err := range errors {
				fmt.Println(err)
			}
		} else {
			fmt.Println("== All Labels are valid")
		}

		fmt.Println("== Results ==")

		maxKeyLength := 0
		for key := range results {
			if len(key) > maxKeyLength {
				maxKeyLength = len(key)
			}
		}
		formatString := fmt.Sprintf("%%-%ds : %%s\n", maxKeyLength)
		for path, value := range results {
			if err != nil {
				log.Errorf("Error renaming path %s key: %v", path, err)
			}
			fmt.Printf(formatString, path, value)
		}
	},
}

func init() {
	extractCmd.Flags().StringVar(&filename, "file", "", "Path to the JSON file")
	extractCmd.Flags().BoolVar(&transformBeforeExtract, "transform", true, "Whether to transform before extracting fields default is true")
	extractCmd.Flags().BoolVar(&rename, "rename", true, "Whether to rename extracted fields")
	err := extractCmd.MarkFlagRequired("file")
	if err != nil {
		log.Fatal(err)
	}
	rootCmd.AddCommand(extractCmd)
}
