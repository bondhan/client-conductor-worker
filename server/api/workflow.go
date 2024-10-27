package api

import (
	"math/rand"
	"net/http"
	"os"

	"github.com/conductor-sdk/conductor-go/sdk/client"
	"github.com/conductor-sdk/conductor-go/sdk/model"
	"github.com/conductor-sdk/conductor-go/sdk/settings"
	"github.com/conductor-sdk/conductor-go/sdk/worker"
	"github.com/conductor-sdk/conductor-go/sdk/workflow/executor"

	log "github.com/sirupsen/logrus"
)

var (
	apiClient = client.NewAPIClient(
		authSettings(),
		httpSettings(),
	)
	taskRunner       = worker.NewTaskRunnerWithApiClient(apiClient)
	workflowExecutor = executor.NewWorkflowExecutor(apiClient)
)

func authSettings() *settings.AuthenticationSettings {
	key := os.Getenv("KEY")
	secret := os.Getenv("SECRET")
	if key != "" && secret != "" {
		return settings.NewAuthenticationSettings(
			key,
			secret,
		)
	}

	return nil
}

func httpSettings() *settings.HttpSettings {
	url := os.Getenv("CONDUCTOR_SERVER_URL")
	if url == "" {
		log.Error("Error: CONDUCTOR_SERVER_URL env variable is not set")
		os.Exit(1)
	}

	return settings.NewHttpSettings(url)
}

func StartWorkflow(w http.ResponseWriter, r *http.Request) {
	min := 5
	max := 10

	// Start the greetings workflow
	id, err := workflowExecutor.StartWorkflow(
		&model.StartWorkflowRequest{
			Name:    "number_square_sleepms",
			Version: 1,
			Input: map[string]int{
				"number": rand.Intn(max-min) + min,
			},
		},
	)

	if err != nil {
		log.Error(err.Error())
		return
	}
	log.Info("Started workflow with Id: ", id)

	// Get a channel to monitor the workflow execution -
	// Note: This is useful in case of short duration workflows that completes in few seconds.
	channel, err := workflowExecutor.MonitorExecution(id)
	if err != nil {
		log.Error(err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{\"error\":" + err.Error() + "}"))
		return
	}
	run := <-channel
	log.Info("Output of the workflow: ", run.Output)

	// Send a success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"id\":" + id + "}"))
}
