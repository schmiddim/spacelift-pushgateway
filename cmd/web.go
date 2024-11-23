package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"spacelift-pushgateway/api"
)

// webCmd represents the web command
var webCmd = &cobra.Command{
	Use:   "web",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		gw := api.NewPushGateway(config.Prometheus.PushGatewayUrl, config.Prometheus.TargetMetric, config.Prometheus.TargetMetricHelp, config.Prometheus.JobName)

		if gw.CheckPushGatewayStatus() != nil {
			log.Error(gw.CheckPushGatewayStatus())
		}

		http.HandleFunc("/health", func(w http.ResponseWriter, request *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
		http.HandleFunc("/push", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
				return
			}
			authHeader := r.Header.Get("Authorization")

			if authHeader != fmt.Sprintf("Bearer %s", apiKey) {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Unable to read request body", http.StatusInternalServerError)
				return
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					log.Fatal("Error closing body")
				}
			}(r.Body)

			// Transform
			for _, splits := range config.Json.ValueSplits {
				body, err = api.TransformJsonValues(body, splits.Path, splits.Separator)
				if err != nil {
					log.Fatalf("Error transforming JSON: %v", err)
				}
			}
			// Extract
			results, err := api.ExtractMultipleJSONPaths(body, config.Json.FieldsToExtract)
			if err != nil {
				log.Fatalf("Error extracting data: %v", err)
			}

			//fmt.Printf("%s", body)
			//fmt.Println(results)
			//Rename
			results, err = api.RenameKeys(results, config.Json.Rename)
			if err != nil {
				log.Error(err)
			}

			//Send to PushGW
			if err := gw.PushMetrics(results); err != nil {
				errMsg := fmt.Sprintf("Failed to push to Pushgateway: %v", err)
				log.Error(errMsg)
				http.Error(w, errMsg, http.StatusInternalServerError)
				return
			}

			log.Info("Successfully pushed data to Pushgateway")

			//fmt.Println(results)
		})
		log.Infof("Server is running on http://localhost:%d", config.App.Port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.App.Port), nil))

	},
}

func init() {
	rootCmd.AddCommand(webCmd)

}
