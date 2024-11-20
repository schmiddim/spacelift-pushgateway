package cmd

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
	"spacelift-pushgateway/api"
	"spacelift-pushgateway/helper"
)

var (
	pushGatewayURL string
	apiKey         string
)
var rootCmd = &cobra.Command{
	Use:   "spacelift-pushgateway",
	Short: "Send Spacelift data to a Prometheus Pushgateway",
	Long: `A lightweight service for forwarding Spacelift job metrics to a Prometheus Pushgateway. 
Requires an API key for authentication and supports configuration via environment variables.`,
	Run: func(cmd *cobra.Command, args []string) {

		gw := api.NewPushGateway(pushGatewayURL)

		if gw.CheckPushGatewayStatus() != nil {
			log.Fatal(gw.CheckPushGatewayStatus())
		}

		log.Infof("Push Gateway: %s", pushGatewayURL)
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

			var payload api.SpaceLiftPayload
			if err := json.Unmarshal(body, &payload); err != nil {
				http.Error(w, "Invalid JSON format", http.StatusBadRequest)
				return
			}

			if err := gw.PushToGateway(payload); err != nil {
				http.Error(w, fmt.Sprintf("Failed to push to Pushgateway: %v", err), http.StatusInternalServerError)
				return
			}

			//log.Info("Successfully pushed data to Pushgateway")
		})
		log.Println("Server is running on :8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
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

	viper.AutomaticEnv()
	pushGatewayURL = viper.GetString("PUSH_GATEWAY_URL")
	apiKey = viper.GetString("API_KEY")
}
