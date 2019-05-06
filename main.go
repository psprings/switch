package main

// import (
// 	"encoding/json"
// 	"errors"
// 	"flag"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"net/url"
// 	"regexp"
// 	"strings"
// 	"time"
// )

import (
	"github.com/psprings/switch/internal/config"
	"github.com/psprings/switch/internal/queues/aws/sqs"
	"github.com/psprings/switch/internal/server"
)

func main() {
	c := config.Retrieve()
	// For testing
	queues := []sqs.PollConfig{
		sqs.PollConfig{
			QueueURL:            c.QueueURL,
			MaxNumberOfMessages: 10,
			VisibilityTimeout:   20,
			WaitTimeSeconds:     0,
			PollInterval:        c.PollInterval,
			BackingURL:          c.BackingURL,
		},
	}
	c.SqsQueues = queues

	if c.Mode == "send" || c.Mode == "both" {
		server.Start(c)
	} else {
		sqs.Receive(queues)
	}
}
