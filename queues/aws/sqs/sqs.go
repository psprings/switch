package sqs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	awsUtils "github.com/psprings/switch/internal/queues/aws"
	"github.com/psprings/switch/internal/utils"
	"github.com/psprings/switch/internal/webhook"
)

// PollConfig : basic configuration for setting up SQS polling
type PollConfig struct {
	QueueURL            string
	MaxNumberOfMessages int64
	WaitTimeSeconds     int64
	VisibilityTimeout   int64
	PollInterval        int
	BackingURL          string
}

// ReceiveMessages :
func (pc *PollConfig) ReceiveMessages(chn chan<- *sqs.Message) {
	log.Printf("Receiving messages for %s", pc.QueueURL)
	sess := awsUtils.InitSession()
	svc := sqs.New(sess)
	output, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            &pc.QueueURL,
		MaxNumberOfMessages: aws.Int64(pc.MaxNumberOfMessages),
		WaitTimeSeconds:     aws.Int64(pc.WaitTimeSeconds),
		VisibilityTimeout:   aws.Int64(pc.VisibilityTimeout),
	})

	if err != nil {
		log.Printf("failed to fetch sqs message %#v", err)
	}

	for _, message := range output.Messages {
		chn <- message
	}
}

// SendMessage :
func (pc *PollConfig) SendMessage(hookURL string, message string) (*sqs.SendMessageOutput, error) {
	sess := awsUtils.InitSession()
	svc := sqs.New(sess)
	result, err := svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"HookURL": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(hookURL),
			},
			"Received": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(utils.CurrentTimeString()),
			},
		},
		MessageBody: aws.String(message),
		QueueUrl:    &pc.QueueURL,
	})

	return result, err
}

// Poll :
func (pc *PollConfig) Poll(chn chan<- *sqs.Message) {
	pollInterval := pc.PollInterval
	pc.ReceiveMessages(chn)
	if pollInterval > 0 {
		for range time.Tick(time.Second * time.Duration(pollInterval)) {
			pc.ReceiveMessages(chn)
		}
	}
}

func basicSQSSvc() *sqs.SQS {
	sess := awsUtils.InitSession()
	svc := sqs.New(sess)
	return svc
}

// HandleMessage :
func HandleMessage(message *sqs.Message) {
	log.Println(*message.MessageId, *message.Body)
}

// DeleteMessage :
func (pc *PollConfig) DeleteMessage(message *sqs.Message) {
	svc := basicSQSSvc()
	_, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &pc.QueueURL,
		ReceiptHandle: message.ReceiptHandle,
	})

	if err != nil {
		fmt.Println("delete error", err)
		return
	}

	log.Println(*message.MessageId, "[DELETED]")
}

// Receive :
func Receive(queues []PollConfig) {
	ReceiveFunc(queues, HandeMessage)
}

// Receive :
func ReceiveFunc(queues []PollConfig, handleMessage func(*sqs.Message)) {
	maxMessages := 10
	chnMessages := make(chan *sqs.Message, maxMessages)

	for _, queue := range queues {
		go queue.Poll(chnMessages)

		log.Printf("Listening on stack queue: %s", queue.QueueURL)

		for message := range chnMessages {
			handleMessage(message)
			queue.DeleteMessage(message)
		}
	}
}

// SendMessageConfig :
type SendMessageConfig struct {
	Queues []PollConfig
}

func matchingQueue(urlToMatch string, queues []PollConfig) (PollConfig, error) {
	for _, queue := range queues {
		if queue.BackingURL == urlToMatch {
			return queue, nil
		}
	}
	return PollConfig{}, errors.New("no matching queue found")
}

// HandleHook :
func HandleHook(w http.ResponseWriter, r *http.Request, queues []PollConfig) {
	postToURL := utils.GetHookURLFromRequest(r)
	log.Printf("match URL: %s", postToURL)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("readAll Body ERROR: %#v", err)
	}
	hook := webhook.Config{
		Header:     r.Header,
		Body:       body,
		Host:       r.Host,
		RemoteAddr: r.RemoteAddr,
		Form:       r.Form,
	}
	postURLErr := utils.TestEndpoint(postToURL)
	if postURLErr != nil {
		log.Println(postToURL)
		queue, err := matchingQueue(postToURL, queues)
		if err != nil {
			log.Printf("matchingQueue error: %#v", err)
		}
		bodyContent := utils.EnsureBodyContent(string(body))
		result, err := queue.SendMessage(postToURL, bodyContent)
		if err != nil {
			log.Printf("sendMessage error: %#v", err)
		}
		hook.MessageID = *result.MessageId
	} else {

	}
	b, err := json.Marshal(hook)
	if err != nil {
		log.Printf("json convert error: %#v", err)
	}
	w.Write(b)
}
