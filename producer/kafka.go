package main

import "github.com/IBM/sarama"

func NewClusterAdmin(brokerAddrs []string) (sarama.ClusterAdmin, error) {
	admin, err := sarama.NewClusterAdmin(brokerAddrs, nil)
	if err != nil {
		return nil, err
	}

	return admin, nil
}

func NewKafkaConfig() *sarama.Config {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V3_6_0_0
	//cfg.Producer.Partitioner = sarama.NewRandomPartitioner
	//cfg.Producer.RequiredAcks = sarama.NoResponse
	cfg.Producer.Return.Successes = true
	return cfg
}

func NewKafkaProducer(brokerAddrs []string, cfg *sarama.Config) (sarama.SyncProducer, error) {
	producer, err := sarama.NewSyncProducer(brokerAddrs, cfg)
	if err != nil {
		return nil, err
	}

	return producer, nil
}
