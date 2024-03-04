package main

import (
	"fmt"
	"github.com/IBM/sarama"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

type config struct {
	brokerAddrs []string
}

type application struct {
	admin    sarama.ClusterAdmin
	producer sarama.SyncProducer
	logger   *slog.Logger
	config   config
}

func main() {
	brokerAddrsEnv := lookupEnvOrDefault("APP_BROKER_ADDRS", "localhost:9094,localhost:9095")
	brokerAddrs := strings.Split(brokerAddrsEnv, ",")

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg := config{
		brokerAddrs: brokerAddrs,
	}

	cl, err := NewClusterAdmin(brokerAddrs)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer cl.Close()

	kafkaConfig := NewKafkaConfig()

	producer, err := NewKafkaProducer(brokerAddrs, kafkaConfig)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer producer.Close()

	app := &application{
		admin:    cl,
		producer: producer,
		logger:   logger,
		config:   cfg,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /produce", app.appHandler(app.produce))

	logger.Info(fmt.Sprintf("server started on :4000"))

	topics, _ := cl.ListTopics()
	for topicName, topic := range topics {
		logger.Info(fmt.Sprintf("Topic name: %s, Topic info: %+v", topicName, topic))
	}

	err = http.ListenAndServe(":4000", mux)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
