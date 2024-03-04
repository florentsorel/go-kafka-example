# Go kafka example

```bash
docker compose up conduktor-console -d
```

Subscribe to the topic `actor` with the consumer
```bash
go run ./consumer
```

Run http server with a `/produce` endpoint to produce a message to kafka.
```bash
go run ./producer
```

Produce message with the endpoint
```bash
curl -X POST -d '{"name": "Bryan Cranston"}' http://localhost:4000/produce
```

The consumer will log the message received from the producer.
