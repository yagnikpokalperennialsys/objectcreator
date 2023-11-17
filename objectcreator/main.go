package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
)

const (
	sqsQueueURL  = "https://sqs.eu-central-1.amazonaws.com/996985152674/article"
	s3BucketName = "objectcreator"
	s3FolderName = "documents/"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
	})
	if err != nil {
		return err
	}

	s3Client := s3.New(sess)
	sqsClient := sqs.New(sess)

	// Receive up to 5 messages from the SQS queue
	result, err := sqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(sqsQueueURL),
		MaxNumberOfMessages: aws.Int64(5),
	})
	if err != nil {
		return err
	}

	for _, message := range result.Messages {
		// Extract the message body from the SQS message
		messageBody := *message.Body

		// Unmarshal the JSON message
		var messageData map[string]string
		if err := json.Unmarshal([]byte(messageBody), &messageData); err != nil {
			log.Printf("Error unmarshaling JSON: %v", err)
			continue
		}

		// Extract UUID and message from JSON
		uuid := messageData["uuid"]
		content := messageData["message"]

		// Create the filename with the path
		filename := s3FolderName + uuid + ".txt"

		// Upload the message content to S3
		err := uploadToS3(s3Client, filename, content)
		if err != nil {
			log.Printf("Error uploading to S3: %v", err)
			continue
		}

		// Delete the message from the SQS queue
		_, err = sqsClient.DeleteMessage(&sqs.DeleteMessageInput{
			QueueUrl:      aws.String(sqsQueueURL),
			ReceiptHandle: message.ReceiptHandle,
		})
		if err != nil {
			log.Printf("Error deleting message from SQS: %v", err)
			continue
		}
	}

	return nil
}

func uploadToS3(s3Client *s3.S3, filename string, content string) error {
	// Create an S3 PutObjectInput
	input := &s3.PutObjectInput{
		Bucket: aws.String(s3BucketName),
		Key:    aws.String(filename),
		Body:   strings.NewReader(content),
	}

	// Upload the file to S3
	_, err := s3Client.PutObject(input)
	if err != nil {
		return err
	}

	fmt.Printf("File '%s' uploaded to S3\n", filename)
	return nil
}
