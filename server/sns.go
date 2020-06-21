package server

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type SubscribeRequest struct {
}

var SNSClient *sns.Client

func InitAWS(wg *sync.WaitGroup) {
	defer wg.Done()
	accessKeyId := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	if accessKeyId == "" || secretKey == "" {
		panic("Missing AWS credentials")
	}

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("Unable to load SDK config, " + err.Error())
	}

	cfg.Region = "us-east-1"
	cfg.Credentials = aws.NewStaticCredentialsProvider(accessKeyId, secretKey, "")

	SNSClient = sns.New(cfg)

	fmt.Println("Successfully initialized AWS SNS Client")
}

func HandleSubscribe(w http.ResponseWriter, r *http.Request) {
	auth, values := AuthenticateRequest(w, r)
	if !auth {
		ErrorUnauthorized(w, r)
		return
	}

	userID := values[KeyUserID].(uint)

	switch r.Method {
	case http.MethodGet:
		fmt.Println("GET /notifications/subscribe")
	case http.MethodPost:
		handleSubscribePost(w, r, userID)
	case http.MethodPut:
		fmt.Println("PUT /notifications/subscribe")
	case http.MethodDelete:
		fmt.Println("DELETE /notifications/subscribe")
	default:
		ErrorMethodNotAllowed(w, r)
	}
}

func handleSubscribePost(w http.ResponseWriter, r *http.Request, userID uint) {
	fmt.Println("POST /notifications/subscribe")

	var user User
	if DB.First(&user, userID).RecordNotFound() {
		fmt.Printf("Failed to find User with ID %d\n", userID)
		ErrorInternalServerError(w)
		return
	}

	if user.SubscriptionSnsStatusId == SNSSubscribedID {
		ErrorBadRequest(w, r, fmt.Errorf("%s", "User already subscribed"))
		return
	}

	if user.MobilePhone == "" {
		ErrorBadRequest(w, r, fmt.Errorf("%s", "User missing MobilePhone"))
		return
	}

	_, errConv := strconv.Atoi(user.MobilePhone)
	if errConv != nil || len(user.MobilePhone) != 10 {
		ErrorBadRequest(w, r, fmt.Errorf("%s", "Invalid User MobilePhone"))
		return
	}

	ctx := context.Background()
	req := SNSClient.SubscribeRequest(&sns.SubscribeInput{
		Endpoint:              aws.String(fmt.Sprintf("+1%s", user.MobilePhone)),
		Protocol:              aws.String("sms"),
		ReturnSubscriptionArn: aws.Bool(true),
		TopicArn:              aws.String("arn:aws:sns:us-east-1:454863778791:Calendays"),
	})

	_, errSend := req.Send(ctx)
	if errSend != nil {
		fmt.Println("SNS subscription failed")
		fmt.Println(errSend)
		ErrorBadGatewayAWS(w)
		return
	}

	user.SubscriptionSnsStatusId = SNSSubscribedID
	if err := DB.Save(&user).Error; err != nil {
		fmt.Printf("Failed to update User (%d) subscription status\n", userID)
		ErrorInternalServerError(w)
		return
	}

	fmt.Println("Successfully subscribed user to SNS topic")
}

func SendSMS(phoneNumber string, message string) error {
	req := SNSClient.PublishRequest(&sns.PublishInput{
		Message:     aws.String(message),
		PhoneNumber: aws.String(fmt.Sprintf("+1%s", phoneNumber)),
	})

	_, err := req.Send(context.Background())
	return err
}
