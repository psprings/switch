package config

import (
	"flag"

	"github.com/psprings/switch/internal/queues/aws/sqs"
)

// Config :
type Config struct {
	Mode       string
	QueueURL   string
	BackingURL string
	SqsQueues  []sqs.PollConfig
}

// Retrieve :
func Retrieve() *Config {
	var mode string
	var sqsQueueURL string
	var hookToProxy string
	flag.StringVar(&hookToProxy, "hook-to-proxy", "http://localhost:8080/github-webhook", "Webhook to proxy")
	flag.StringVar(&sqsQueueURL, "sqs-queue-url", "", "SQS Queue URL for testing")
	flag.StringVar(&mode, "mode", "receive", "Can be one of [send, receive]. Send will forward webhooks to queue, receive will process and POST.")
	flag.Parse()
	return &Config{
		Mode:       mode,
		QueueURL:   sqsQueueURL,
		BackingURL: hookToProxy,
	}
}
