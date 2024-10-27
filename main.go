package main

import (
	"client-conductor-worker/src"
	"fmt"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"

	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"

	"github.com/conductor-sdk/conductor-go/sdk/client"
	"github.com/conductor-sdk/conductor-go/sdk/settings"

	"github.com/conductor-sdk/conductor-go/sdk/worker"
	"github.com/conductor-sdk/conductor-go/sdk/workflow/executor"
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
	log.Println("conductor server:", url)
	if url == "" {
		log.Error("Error: CONDUCTOR_SERVER_URL env variable is not set")
		os.Exit(1)
	}

	return settings.NewHttpSettings(url)
}

func main() {

	cpu, err := strconv.Atoi(os.Getenv("GOMAXPROCS"))
	if err != nil {
		log.Fatalf("Error: GOMAXPROCS env variable is not set %s", err)
	}
	runtime.GOMAXPROCS(cpu)

	log.Println("GOMAXPROCS is set to:", runtime.GOMAXPROCS(0))

	batchSizeS := os.Getenv("BATCH_SIZE")
	if strings.TrimSpace(batchSizeS) == "" {
		log.Fatal("Error: BATCH_SIZE env variable is not set")
	}
	batchSize, err := strconv.Atoi(batchSizeS)
	if err != nil {
		log.Fatalf("Error: BATCH_SIZE env variable is not set")
	}
	log.Println("BATCH_SIZE is set to:", batchSizeS)

	pollingTimeS := os.Getenv("POLLING_TIME")
	if strings.TrimSpace(pollingTimeS) == "" {
		log.Fatal("Error: POLLING_TIME env variable is not set")
	}
	pollingTime, err := strconv.Atoi(pollingTimeS)
	if err != nil {
		log.Fatalf("Error: POLLING_TIME env variable is not a number: %v", err)
	}
	log.Println("POLLING_TIME is set to:", pollingTimeS)

	err = taskRunner.StartWorker("number", conductorworker.Number, batchSize, time.Duration(pollingTime)*time.Millisecond)
	if err != nil {
		log.Fatalf("error starting worker, err: %s", err)
	}

	err = taskRunner.StartWorker("square", conductorworker.Square, batchSize, time.Duration(pollingTime)*time.Millisecond)
	if err != nil {
		log.Fatalf("error starting worker, err: %s", err)
	}

	err = taskRunner.StartWorker("sleepms", conductorworker.Sleepms, batchSize, time.Duration(pollingTime)*time.Millisecond)
	if err != nil {
		log.Fatalf("error starting worker, err: %s", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Blocking, press ctrl+c to continue...")
	<-done // Will block here until user hits ctrl+c
}
