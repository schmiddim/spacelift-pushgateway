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
	viper.AutomaticEnv()
	apiKey = viper.GetString("API_KEY")

}
