package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"spacelift-pushgateway/api"
	"spacelift-pushgateway/helper"
)

func readJsonFile(filePath string) []byte {
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	return jsonData
}

type ValueSplits struct {
	Path      string
	Separator string
}

type Config struct {
	App struct {
		Port int
	}

	Json struct {
		ValueSplits     []ValueSplits
		FieldsToExtract []string
		Rename          []api.Rename
	}
	Logging struct {
		Level  string
		Format string
	}
	Prometheus struct {
		PushGatewayUrl   string
		TargetMetric     string
		TargetMetricHelp string
		JobName          string
	}
}

var (
	apiKey string
	config Config
)
var rootCmd = &cobra.Command{
	Use:   "spacelift-pushgateway",
	Short: "Send Spacelift data to a Prometheus Pushgateway",
	Long: `A lightweight service for forwarding Spacelift job metrics to a Prometheus Pushgateway. 
Requires an API key for authentication and supports configuration via environment variables.`,
	Run: func(cmd *cobra.Command, args []string) {

		//gw := api.NewPushGateway(config.Prometheus.PushGatewayUrl, config.Prometheus.TargetMetric, config.Prometheus.FieldsToExtract)
		//
		//if gw.CheckPushGatewayStatus() != nil {
		//	log.Fatal(gw.CheckPushGatewayStatus())
		//}
		//
		//log.Infof("Push Gateway: %s", config.Prometheus.PushGatewayUrl)
		//

		//http.HandleFunc("/push", func(w http.ResponseWriter, r *http.Request) {
		//	if r.Method != http.MethodPost {
		//		http.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
		//		return
		//	}
		//	authHeader := r.Header.Get("Authorization")
		//
		//	if authHeader != fmt.Sprintf("Bearer %s", apiKey) {
		//		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		//		return
		//	}
		//	body, err := io.ReadAll(r.Body)
		//	if err != nil {
		//		http.Error(w, "Unable to read request body", http.StatusInternalServerError)
		//		return
		//	}
		//	defer func(Body io.ReadCloser) {
		//		err := Body.Close()
		//		if err != nil {
		//			log.Fatal("Error closing body")
		//		}
		//	}(r.Body)
		//
		//	var payload api.SpaceLiftPayload
		//	if err := json.Unmarshal(body, &payload); err != nil {
		//		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		//		return
		//	}
		//	fieldsToExtract := []string{"branch", "name", "namespace", "commit.author", "commit.hash", "labels"}
		//	extractor := api.NewExtractor(fieldsToExtract)
		//	extractedData, err := extractor.Extract(body)
		//
		//	fmt.Println(extractedData)
		//	if err := gw.PushToGateway(payload); err != nil {
		//		http.Error(w, fmt.Sprintf("Failed to push to Pushgateway: %v", err), http.StatusInternalServerError)
		//		return
		//	}

		//log.Info("Successfully pushed data to Pushgateway")
		//})
		log.Infof("Server is running on :%d", config.App.Port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.App.Port), nil))
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	helper.LoggerInit()
	viper.SetDefault("PUSH_GATEWAY_URL", "http://localhost:9091")
	viper.SetDefault("API_KEY", "extreme-secret-key")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error Loading config: %v", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Fehler unmarshalling config: %v", err)
	}
	log.Info(config)
	viper.AutomaticEnv()
	apiKey = viper.GetString("API_KEY")

}
