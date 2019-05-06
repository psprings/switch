package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/psprings/switch/internal/config"
	"github.com/psprings/switch/internal/queues/aws/sqs"
)

// Start :
func Start(c *config.Config) {
	port := 8000
	listenOn := fmt.Sprintf("0.0.0.0:%d", port)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Will abstract behind generic Handler Func later to support multiple queue types
		sqs.HandleHook(w, r, c.SqsQueues)
	})
	if err := http.ListenAndServe(listenOn, nil); err != nil {
		panic(err)
	}
	log.Printf("Listening on %s", listenOn)
}
