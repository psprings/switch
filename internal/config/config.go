package config

import (
	"flag"

	"github.com/psprings/switch/internal/queues/aws/sqs"
)

// Config :
type Config struct {
	Mode         string
	QueueURL     string
	BackingURL   string
	PollInterval int
	SqsQueues    []sqs.PollConfig
}

// Retrieve :
func Retrieve() *Config {
	var mode string
	var sqsQueueURL string
	var hookToProxy string
	var pollInterval int
	flag.StringVar(&hookToProxy, "hook-to-proxy", "http://localhost:8080/github-webhook", "Webhook to proxy")
	flag.StringVar(&sqsQueueURL, "sqs-queue-url", "", "SQS Queue URL for testing")
	flag.IntVar(&pollInterval, "interval", 15, "Interval in seconds between queue checks")
	flag.StringVar(&mode, "mode", "receive", "Can be one of [send, receive]. Send will forward webhooks to queue, receive will process and POST.")
	flag.Parse()
	return &Config{
		Mode:       mode,
		QueueURL:   sqsQueueURL,
		BackingURL: hookToProxy,
	}
}
