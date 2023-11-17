package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sns"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io/ioutil"
)

var ginLambda *ginadapter.GinLambda

func init() {
	r := gin.Default()
	r.POST("/", PostDocument)
	r.GET("/:id", GetDocument)
	ginLambda = ginadapter.New(r)
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if ginLambda == nil {
		panic("ginLambda should not be nil")
	}

	// Pass the API Gateway proxy request directly to the Gin Lambda adapter
	response, err := ginLambda.ProxyWithContext(ctx, req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, err
	}

	return response, nil
}

func main() {
	lambda.Start(handler)
}

func GetDocument(c *gin.Context) {
	documentID := c.Param("id")

	// Placeholder for fetching document content from S3 based on documentID
	documentContent, err := fetchDocumentFromS3(documentID)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	c.String(200, documentContent)
}

func PostDocument(c *gin.Context) {

	// Read the request body into a byte slice
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to read request body"})
		return
	}

	// Generate UUID
	newUUID := uuid.New()

	// Placeholder for sending JSON to SNS topic
	err = sendTextToSNSTopic(newUUID.String(), string(body))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Send UUID as response along with status code
	c.JSON(200, gin.H{"id": newUUID.String(), "body": string(body)})
}

func fetchDocumentFromS3(documentID string) (string, error) {
	// Create a new AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
	})
	if err != nil {
		return "", err
	}
	fmt.Println("Session connected")

	// Create an S3 service client
	s3Client := s3.New(sess)

	// Specify your S3 bucket and key based on the documentID
	bucketName := "objectcreator"
	s3Key := "documents/" + documentID + ".txt"
	fmt.Println("S3 key is", s3Key)

	// Get the object from S3
	output, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(s3Key),
	})
	if err != nil {
		return "", err
	}
	defer output.Body.Close() // Ensure the response body is closed after use

	// Read the content from the S3 object
	content, err := ioutil.ReadAll(output.Body)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func sendTextToSNSTopic(uuid string, text string) error {
	// Create a new AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
	})
	if err != nil {
		return err
	}

	// Create an SNS service client
	snsClient := sns.New(sess)

	// Specify your SNS topic ARN
	topicArn := "arn:aws:sns:eu-central-1:996985152674:mytopic"

	// Create a map to represent the JSON message
	messageData := map[string]interface{}{
		"uuid":    uuid,
		"message": text,
	}

	// Convert the map to JSON
	jsonMessage, err := json.Marshal(messageData)
	if err != nil {
		return err
	}

	// Publish the JSON message to the SNS topic
	_, err = snsClient.Publish(&sns.PublishInput{
		TopicArn: aws.String(topicArn),
		Message:  aws.String(string(jsonMessage)),
		MessageAttributes: map[string]*sns.MessageAttributeValue{
			"eventType": {
				DataType:    aws.String("String"),
				StringValue: aws.String("sqs_specific_event"),
			},
		},
	})

	if err != nil {
		return err
	}

	fmt.Printf("Message with UUID %s published to SNS\n", uuid)
	return nil
}
